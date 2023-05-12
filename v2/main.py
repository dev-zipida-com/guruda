# """Return docs selected using the maximal marginal relevance.

# Maximal marginal relevance optimizes for similarity to query AND diversity
# among selected documents.

# Args:
#     query: Text to look up documents similar to.
#     k: Number of Documents to return. Defaults to 4.
#     fetch_k: Number of Documents to fetch to pass to MMR algorithm.
#     lambda_mult: Number between 0 and 1 that determines the degree
#                 of diversity among the results with 0 corresponding
#                 to maximum diversity and 1 to minimum diversity.
#                 Defaults to 0.5.
# Returns:
#     List of Documents selected by maximal marginal relevance.
# """

from langchain.embeddings.openai import OpenAIEmbeddings
from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain.vectorstores import DeepLake
from langchain.document_loaders import TextLoader
from langchain.chat_models import ChatOpenAI
from langchain import SerpAPIWrapper
from langchain.agents import Tool, AgentExecutor, LLMSingleActionAgent, AgentOutputParser
from langchain.memory import ConversationBufferMemory
from langchain.agents import initialize_agent
from langchain.agents import AgentType
from dotenv import load_dotenv
from github import Github
import deeplake
import os

load_dotenv()
embeddings = OpenAIEmbeddings(openai_api_key=os.getenv("OPENAI_API_KEY"),disallowed_special=())

def get_github_contents(addr, branch_name):

    contents_list = []
    access_token = os.getenv("GITHUB_API_KEY")
    g = Github(access_token)
    repo = g.get_repo(addr)

    contents = repo.get_contents('', ref=branch_name)
    contents_list = []

    text_splitter = RecursiveCharacterTextSplitter(
        chunk_size=1000,
        chunk_overlap=100,
        length_function = len,
    )
    # GET github repository contents using Github API
    while contents:
        file_content = contents.pop(0)
        extensions = (".go", ".py", ".js", ".ts", ".tsx", ".html", ".css", ".md", ".java", ".c", ".cpp")

        if file_content.type == 'dir':
            contents.extend(repo.get_contents(file_content.path, ref=branch_name))
        else:
            file_extension = os.path.splitext(file_content.path)[1]
            if file_extension not in extensions:
                continue
            contents_list.extend(text_splitter.create_documents([file_content.decoded_content.decode("utf-8")]))

    return contents_list    


# Get Repository from Github API response
user_name = input("Enter your github username: ")
repo_name = input("Enter the repo name: ")
branch_name = input("Enter the branch name: ")
addr = user_name + '/' + repo_name

username = os.getenv("DEEPLAKE_USERNAME")
dataset_path = f"hub://{username}/{user_name}_{repo_name}"

# DeepLake Vector Database connect
if deeplake.exists(path=dataset_path):
    db = DeepLake(token=os.getenv("ACTIVELOOP_TOKEN"), dataset_path=dataset_path, embedding_function=embeddings)

else:
    db = DeepLake(token=os.getenv("ACTIVELOOP_TOKEN"), dataset_path=dataset_path, embedding_function=embeddings, public=True)

    contents_list = get_github_contents(addr=addr, branch_name=branch_name)

    # Add repository contents in Deeplake DB
    db.add_documents(contents_list)

# database search option
retriever = db.as_retriever()
retriever.search_kwargs['distance_metric'] = 'cos' 
retriever.search_kwargs['fetch_k'] = 100
retriever.search_kwargs['maximal_marginal_relevance'] = True
retriever.search_kwargs['k'] = 10
# retriever.search_kwargs['lambda_mult'] = 0.5

# using GPT 4 API
# model = ChatOpenAI(model_name='gpt-4', openai_api_key=os.getenv("OPENAI_KEY_GPT4"))

# using GPT 3.5 API
model = ChatOpenAI(model_name='gpt-3.5-turbo', openai_api_key=os.getenv("OPENAI_API_KEY"))

# If you want to use the search API, use this option.
search = SerpAPIWrapper(serpapi_api_key=os.getenv("SERPAPI_API_KEY"))

search_tool = [
    Tool(
        name = "Searching with VectorStore",
        func=retriever.get_relevant_documents,
        description="Search similar documents stored in the VectorStore: DeepLake"
    ),
    # If you want to use the search API, use this option.
    Tool(
        name = "Search by SerpAPI",
        func=search.run,
        description="this is a tool for searching web by SerpAPI. check Vector DB first, and if The data stored in Vector DB is not sufficient, use this tool."
    ),
]

# Create Agent
memory = ConversationBufferMemory(memory_key="chat_history", return_messages=True)
agent_chain = initialize_agent(tools=search_tool, llm=model, 
                               agent=AgentType.CHAT_CONVERSATIONAL_REACT_DESCRIPTION, 
                               verbose=True, memory=memory)

while True:
    question = input("question: ")

    if question == "exit":
        break

    result = (agent_chain.run(input=question))

    print(f"-> **Question**: {question} \n")
    print(f"**Answer**: {result} \n")
    print("-------------------------")
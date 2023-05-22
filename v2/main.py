from langchain.embeddings.openai import OpenAIEmbeddings
from langchain.vectorstores import DeepLake
from langchain.chat_models import ChatOpenAI
from langchain import SerpAPIWrapper
from langchain.agents import Tool
from langchain.memory import ConversationBufferMemory
from langchain.agents import initialize_agent
from langchain.agents import AgentType
from get_content import get_github_contents
from dotenv import load_dotenv
import deeplake
import os

load_dotenv()
embeddings = OpenAIEmbeddings(openai_api_key=os.getenv("OPENAI_API_KEY"),disallowed_special=())

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
model = ChatOpenAI(model_name='gpt-4', openai_api_key=os.getenv("OPENAI_KEY_GPT4"))

# using GPT 3.5 API
# model = ChatOpenAI(model_name='gpt-3.5-turbo', openai_api_key=os.getenv("OPENAI_API_KEY"))

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
# guruda

전체 코드의 맥락을 인식하는 Automatic ReadMe Generator &amp; Refactoring machine.

<br>

## 1. Setup Environment

```
$ pip install -f requirements.txt
```

<br>

## 2. Environment Variables

v2 폴더내에 .env 파일을 생성한다.

- OPENAI_API_KEY
  - [OPENAI API keys](https://platform.openai.com/account/api-keys)
- GITHUB_API_KEY
  - [GITHUB API Token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)
- DEEPLAKE_USERNAME
  - [DEEPLAKE](https://www.deeplake.ai/) 회원 username
- ACTIVELOOP_TOKEN
  - https://app.activeloop.ai/profile/{username}/apitoken

<br>

## 3. Run `main.py`

```
$ cd v2
$ python3 main.py
```

<br>

## 4. Enter User Information

- Example

```bash
Enter your github username: (dev-zipida-com)
Enter the repo name: (spo-vdvs-system)
Enter the branch name: (develop)

...

question: (please explain generate-link modules made by golang. refer to the Vector Store)
```

- 현재 프롬프트 미완성으로 인해 `refer to the Vector Store` 라는 문구를 질문 마지막에 붙여줘야 repository 의 코드 정보를 참고하여 답변을 받을 수 있습니다.

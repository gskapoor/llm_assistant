## Set up environment

Create a virtual environment

```sh
python -m venv .venv
source .venv/bin/activate
which pip # make sure it output the pip points to your virtual environment
pip install -r requirements.txt
```

## Start the server

```sh
uvicorn main:app --reload
```


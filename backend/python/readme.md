## Set up environment

Create a virtual environment

```sh
python -m venv .venv
source .venv/bin/activate # import virtual environment
which pip # make sure it output the pip points to your virtual environment
pip install -r requirements.txt
```

On Windows, please use pwsh 7+, and install python, then

```ps1
python -m venv .venv
& .venv/Scripts/Activate.ps1 # import virtual environment
where.exe pip # make sure it output the pip points to your virtual environment
pip install -r requirements.txt
```

## Start the server

Make sure import the virtual environment and then start the server with

```sh
uvicorn main:app --reload
```


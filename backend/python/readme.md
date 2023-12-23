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

## Read doc and try API

Open the link `http://localhost:8000/docs` in the browser and read the doc try the API by clicking "Try it out". 

No need curl or httpie, yay!


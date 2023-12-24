from assistant import initialize_conversation, continue_conversation, end_conversation, AssistantSession, AssistantSessionMessage
from fastapi import FastAPI

app = FastAPI()

# Implements the functions in assistant.py as an API
@app.get("/")
async def read_root():
    return {"Hello": "World"}

@app.get("/assistant")
async def init():
    new_assistant = await initialize_conversation()
    return {"assistant_session": new_assistant}

@app.post("/assistant")
async def cont(assistant_message: AssistantSessionMessage):
    response = await continue_conversation(assistant_message)
    return {"response": response}

@app.delete("/assistant")
async def end(assistant: AssistantSession):
    deleted_status = await end_conversation(assistant)
    return {"status": deleted_status}

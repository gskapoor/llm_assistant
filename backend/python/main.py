from low_level_assistant import initialize_conversation, continue_conversation, end_conversation, AssistantSession
from fastapi import FastAPI

app = FastAPI()

# test session
@app.get("/")
async def read_root():
    return {"Hello": "World"}

@app.post("/assistant")
async def init():
    new_assistant = initialize_conversation()
    return {"assistant_session": new_assistant}

@app.get("/assistant")
async def cont(assistant: AssistantSession, message: str):
    response = continue_conversation(assistant, message)
    return {"response": response}

@app.delete("/assistant")
async def end(assistant: AssistantSession):
    deleted_status = end_conversation(assistant)
    return {"status": deleted_status}
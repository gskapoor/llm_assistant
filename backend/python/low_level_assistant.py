import asyncio
import time
from dotenv import load_dotenv
from openai import AsyncOpenAI
from pydantic import BaseModel

class AssistantSession(BaseModel):
  """
  """
  assistant_id: str
  thread_id: str

WAIT_TIME = 3 # (seconds) the time to wait between checking whether the LLM's response is complete
MAX_NUM_WAITS = 1000 # sets wait time for timeout from waiting on model response

load_dotenv()
client = AsyncOpenAI()

async def initialize_conversation():
  """TODO: write doc here
  """
  assistant = await client.beta.assistants.create(
    name="Personal Assistant",
    instructions="You are a personal assistant whose role is to help the user complete their tasks and to entertaing the user via conversation.",
    model="gpt-3.5-turbo-1106",
  )

  thread = await client.beta.threads.create()

  new_assistant = AssistantSession(assistant_id=assistant.id,thread_id=thread.id)

  return new_assistant


async def continue_conversation(current_assistant, next_message):
  threadId = current_assistant.thread_id
  assistantId = current_assistant.assistant_id

  await client.beta.threads.messages.create(
    thread_id=threadId, role="user", content=next_message
  )

  run = await client.beta.threads.runs.create(
    thread_id=threadId, assistant_id=assistantId
  )

  num_waits = 0

  while num_waits < MAX_NUM_WAITS:
    run = await client.beta.threads.runs.retrieve(
      thread_id=threadId, run_id=run.id
    )

    match run.status:
      case "requires_action":
        continue # This means that a tool/function needs to be called. Will implement once we have at least one tool.
      case "cancelled":
        print("cancelled")
        break
      case "failed":
        print("failed")
        break
      case "expired":
        print("expired")
        break
      case "completed":
        break
      case _:
        num_waits += 1
        time.sleep(WAIT_TIME)

  messages = await client.beta.threads.messages.list(thread_id=threadId)
  response = messages.data[0].content[0].text.value

  return response


async def end_conversation(assistant_to_delete):
  thread_deleted = await client.beta.threads.delete(assistant_to_delete.thread_id)
  assistant_deleted = await client.beta.assistants.delete(assistant_to_delete.assistant_id)

  return thread_deleted.deleted and assistant_deleted.deleted

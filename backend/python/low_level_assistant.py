from dotenv import load_dotenv
import asyncio
import time

load_dotenv()
from openai import AsyncOpenAI

WAIT_TIME = 3 # the time to wait between checking whether the LLM's response is complete
MAX_NUM_WAITS = 1000 # sets wait time for timeout from waiting on model response
client = AsyncOpenAI()


async def initialize_conversation():
  assistant = await client.beta.assistants.create(
    name="Personal Assistant",
    instructions="You are a personal assistant whose role is to help the user complete their tasks and to entertaing the user via conversation.",
    model="gpt-3.5-turbo-1106",
  )

  thread = await client.beta.threads.create()

  return assistant, thread


async def continue_conversation(assistant, thread, next_message):
  message = await client.beta.threads.messages.create(
    thread_id=thread.id, role="user", content=next_message
  )

  run = await client.beta.threads.runs.create(
    thread_id=thread.id, assistant_id=assistant.id
  )

  num_waits = 0

  while num_waits < MAX_NUM_WAITS:
    run = await client.beta.threads.runs.retrieve(
      thread_id=thread.id, run_id=run.id
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

  messages = await client.beta.threads.messages.list(thread_id=thread.id)

  return messages.data[0].content[0].text.value


async def start_conversation(initial_message):
  assistant, thread = await initialize_conversation()

  messages = await continue_conversation(assistant, thread, initial_message)

  return assistant, thread, messages

async def end_conversation(assistant, thread):
  thread_deleted = await client.beta.threads.delete(thread.id)
  assistant_deleted = await client.beta.assistants.delete(assistant_id=assistant.id)
  return thread_deleted.deleted and assistant_deleted.deleted
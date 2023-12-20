from dotenv import load_dotenv
import asyncio
import time

load_dotenv()
from openai import AsyncOpenAI

wait_time = (
    3  # the time to wait between checking whether the LLM's response is complete
)
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

    while num_waits < 100:
        run = await client.beta.threads.runs.retrieve(
            thread_id=thread.id, run_id=run.id
        )

        match run.status:
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
                time.sleep(wait_time)

    messages = await client.beta.threads.messages.list(thread_id=thread.id)

    return messages.data[0].content[0].text.value


async def start_conversation(initial_message):
    assistant, thread = await initialize_conversation()

    messages = await continue_conversation(assistant, thread, initial_message)

    return assistant, thread, messages


async def main() -> None:
    assistant, thread, my_conversation = await start_conversation("Hi! What is 2 + 2?")
    myer_conversation = await continue_conversation(
        assistant,
        thread,
        "What are the colors of the rainbow? Nothing else after that.",
    )
    print(myer_conversation)


asyncio.run(main())

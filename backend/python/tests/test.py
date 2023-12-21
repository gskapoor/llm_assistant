import unittest
import asyncio
from ..low_level_assistant import start_conversation, continue_conversation, end_conversation

class TestAssistant(unittest.IsolatedAsyncioTestCase):

    async def test_basic(self):
        assistant, thread, my_conversation = await start_conversation("I need you to help me test a piece of software I'm writing. To do so, please respond saying exactly the phrase in quotes but without the quotes: 'I am Groot.' Do not say anything else or acknowledge the request.")

        myer_conversation = await continue_conversation(
            assistant,
            thread,
            "Also repeat the following phrase in the same way: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.'",
        )
        end_worked = await end_conversation(assistant, thread)

        self.assertEqual(my_conversation,"I am Groot.")
        self.assertEqual(myer_conversation, "Lorem ipsum dolor sit amet, consectetur adipiscing elit.")
        self.assertTrue(end_worked)

if __name__ == '__main__':
    unittest.main()
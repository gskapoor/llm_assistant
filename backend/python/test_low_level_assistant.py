import unittest
import asyncio
from assistant import initialize_conversation, continue_conversation, end_conversation, AssistantSession

class TestAssistant(unittest.IsolatedAsyncioTestCase):

    async def test_basic(self):
        assistant = await initialize_conversation()
        my_conversation = await continue_conversation(
            assistant,
            "I need you to help me test a piece of software I'm writing. To do so, please respond saying exactly the phrase in quotes but without any quotes: 'I am Groot.' Do not say anything else or acknowledge the request.",
        )

        myer_conversation = await continue_conversation(
            assistant,
            "Also repeat the following phrase in the same way: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.'",
        )
        end_worked = await end_conversation(assistant)

        self.assertEqual(my_conversation,"I am Groot.")
        self.assertEqual(myer_conversation, "Lorem ipsum dolor sit amet, consectetur adipiscing elit.")
        self.assertTrue(end_worked)

if __name__ == '__main__':
    unittest.main()
import unittest
import asyncio
import requests

class TestAssistantApi(unittest.IsolatedAsyncioTestCase):

    async def full_run_test():
        r = requests.put()
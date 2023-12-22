import { useState } from 'react';
import { AudioRecorder, useAudioRecorder } from 'react-audio-voice-recorder';
import './App.css';

const initialMessages = [{author: "maya", text: "Welcome to MAYA! How may I assist you today?"}];

function App() {
  const [messages, setMessages] = useState(initialMessages);
  const recorderControls = useAudioRecorder();

  async function handleAudioMessage(blob) {
    const url = URL.createObjectURL(blob);
    const audioFile = new File([blob], 'audio.mp3', { type: 'audio/mpeg' });
    const data = new FormData();
    data.append("audio", audioFile);
    const options = {
      method: "POST",
      mode: "cors",
      body: data
    };
    const response = await fetch('http://localhost:8080/voice', options);
    const responseText = await response.text();
    genMessage({author: "user", text: responseText, audio: url});
  };

  async function handleTextMessage(e) {
    e.preventDefault();

    const form = e.target;
    const formData = new FormData(form);
    const formJson = Object.fromEntries(formData.entries());

    genMessage({author: "user", text: formJson.message, audio: null})
    const responseJson = {message: "[insert API response here]"};
    genMessage({author: "maya", text: responseJson.message, audio: null})

  }

  function genMessage(data) {
    setMessages(messages => [{author: data.author, text: data.text, audio: data.audio}, ...messages]);
  }

  function MessageList({messages}) {
    return (
      <div id="messageList">
        {messages.map(message => (
          <div className={"message-" + message.author}>
            <div className="message-box">
              <div className="message-info">
                {message.author}
              </div>
              {message.text != null &&
                <div className="message-text">
                  {message.text}
                </div>
              }
              {message.audio != null &&
                <div className="message-audio">
                  <audio src={message.audio} controls={true}></audio>
                </div>
              }
            </div>
          </div>
        ))}
      </div>
    );
  }

  return (
    <div className="App">
      <div id="chat">
        <MessageList messages={messages} />
        <div id="messageInput">
          <form id="messageInputArea" onSubmit={handleTextMessage}>
            <div className="inputLeft">
              <textarea name="message" type="text" id="messageType" placeholder="Enter message"></textarea>
            </div>
            <div className="inputRight">
              <AudioRecorder 
              onRecordingComplete={(blob) => handleAudioMessage(blob)}
              recorderControls={recorderControls}
              downloadFileExtension="mp3"
              />
              <button type="submit" id="textSend">Send</button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}

export default App;

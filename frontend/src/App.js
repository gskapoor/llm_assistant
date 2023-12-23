import { useState, useRef } from 'react';
import { AudioRecorder, useAudioRecorder } from 'react-audio-voice-recorder';
import './App.css';

const initialMessages = [{author: "maya", text: "Welcome to MAYA! How may I assist you today?"}];

function App() {
  const [messages, setMessages] = useState(initialMessages);
  const [theme, setTheme] = useState("grayscale");
  const recorderControls = useAudioRecorder();
  const [modal, setModal] = useState(null);
  const formRef = useRef(null);

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
    formRef.current.reset();
    const responseJson = {message: "[insert API response here]"};
    genMessage({author: "maya", text: responseJson.message, audio: null})

  }

  function handleEnter(e) {
    if(e.keyCode === 13 && e.shiftKey === false) {
      e.preventDefault();
      formRef.current.requestSubmit();
    }
  }

  function genMessage(data) {
    setMessages(messages => [{author: data.author, text: data.text, audio: data.audio}, ...messages]);
  }

  function toggleModal(mode) {
    if (mode === modal && modal != null) {
      setModal(null)
    }
    else {
      setModal(mode);
    }
  };

  function MessageList({messages}) {
    return (
      <div id="messageList">
        {messages.map(message => (
          <div className={"message " + message.author}>
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

  function Navbar() {
    return (
      <nav id="navId">
        <button onClick={() => toggleModal("account")}>Account</button>
        <button onClick={() => toggleModal("theme")}>Theme</button>
      </nav>
    );
  };

  function Modal() {
    return (
      <div className="modal">
          <div onClick={() => toggleModal(null)} className="overlay"></div>
          {modal === "theme" &&
          <div className="modal-content"> 
            <h2>Theme</h2>
            <select value={theme} onChange={e => setTheme(e.target.value)}>
              <option value="grayscale">Grayscale</option>
              <option value="kawaii">Kawaii</option>
            </select>
            <button className="close-modal" onClick={() => toggleModal(null)}>CLOSE</button>
          </div>
          }
          {modal === "account" &&
          <div className="modal-content"> 
            <h2>Account</h2>
            <p>
            Lorem ipsum dolor sit amet consectetur adipisicing elit. Provident
            perferendis suscipit officia recusandae, eveniet quaerat assumenda
            id fugit, dignissimos maxime non natus placeat illo iusto!
            Sapiente dolorum id maiores dolores? Illum pariatur possimus
            quaerat ipsum quos molestiae rem aspernatur dicta tenetur. Sunt
            placeat tempora vitae enim incidunt porro fuga ea.
            </p>
            <button className="close-modal" onClick={() => toggleModal(null)}>CLOSE</button>
          </div>
          }
      </div>
    );
  }

  return (
    <div className={"App " + theme}>
      {modal && 
        <Modal />
      }
      <div id="chat">
        <MessageList messages={messages} />
        <div id="messageInput">
          <form id="messageInputArea" onSubmit={handleTextMessage} ref={formRef}>
            <div className="inputLeft">
              <textarea name="message" type="text" id="messageType" placeholder="Enter message" onKeyDown={handleEnter}></textarea>
            </div>
            <div className="inputRight">
              <AudioRecorder 
              onRecordingComplete={(blob) => handleAudioMessage(blob)}
              recorderControls={recorderControls}
              />
              <button type="submit" id="textSend">Send</button>
            </div>
          </form>
        </div>
      </div>
      <Navbar />
    </div>
  );
}

export default App;

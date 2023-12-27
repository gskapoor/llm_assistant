import { useState, useRef, useEffect } from 'react';
import { AudioRecorder, useAudioRecorder } from 'react-audio-voice-recorder';
import { googleLogout, useGoogleLogin } from '@react-oauth/google';
import './App.css';

const initialMessages = [{author: "maya", text: "Welcome to MAYA! How may I assist you today?"}];
const API_URL = "http://localhost:8080"

function App() {
  const [messages, setMessages] = useState(initialMessages);
  const [theme, setTheme] = useState("grayscale");
  const [modal, setModal] = useState(null);
  const [ user, setUser ] = useState(null);
  const [ profile, setProfile ] = useState(null);
  const recorderControls = useAudioRecorder();
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
    fetch(API_URL + '/voice', options)
    .then(response => response.json())
    .then(responseJson => {
      genMessage({author: "user", text: responseJson.transcribed_text, audio: url});
      genMessage({author: "maya", text: responseJson.response, audio: null});
    })
  };

  async function handleTextMessage(e) {
    e.preventDefault();

    const form = e.target;
    const formData = new FormData(form);
    const formJson = Object.fromEntries(formData.entries());
    genMessage({author: "user", text: formJson.message, audio: null});

    const options = {
      method: "POST",
      body: formJson.message
    };
    fetch(API_URL + '/text', options)
    .then(response => response.text())
    .then(responseText => {
      genMessage({author: "maya", text: responseText, audio: null});
    })
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

  const login = useGoogleLogin({
    onSuccess: (codeResponse) => setUser(codeResponse),
    onError: (error) => console.log('Login Failed:', error)
  });

  const logout = () => {
      googleLogout();
      setProfile(null);
  };

  useEffect(
    () => {
        if (user) {
          const options = {
            method: "GET",
            headers: {
              Authorization: `Bearer ${user.access_token}`,
              Accept: 'application/json'
            }
          };
          fetch(`https://www.googleapis.com/oauth2/v1/userinfo?access_token=${user.access_token}`, options)
          .then(res => res.json())
          .then(json => setProfile(json))
          .catch(err => console.log(err));
        }
    }, [user]
);

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
            {profile ? (
                <div>
                    <img src={profile.picture} alt="user" />
                    <h3>User Logged in</h3>
                    <p>Name: {profile.name}</p>
                    <p>Email Address: {profile.email}</p>
                    <br />
                    <br />
                    <button onClick={logout}>Log out</button>
                </div>
            ) : (
                <button onClick={() => login()}>Sign in with Google 🚀 </button>
            )}
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
              downloadFileExtension="mp3"
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

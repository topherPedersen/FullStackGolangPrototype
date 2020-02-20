import React from 'react';
import PlayList from './PlayList';

const mainDivStyle = {
  textAlign: "center",
  marginTop: "10em",
};

const soundcloudIframeStyle = {
  width: 600,
  height: 125,
  scrolling: "no",
  frameborder: "no"
};

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      songs: []
    }
  }

  componentDidMount() {
    this.pingServer();
  }

  pingServer() {
    var httpRequest = new XMLHttpRequest();
    var httpRequestURL = "https://topherstop100.ngrok.io";
    httpRequest.open("GET", httpRequestURL, true);
    httpRequest.setRequestHeader('Content-Type', 'application/json');
    var setStateAsync = this.setState.bind(this)
    httpRequest.onreadystatechange = function() {
      if (this.readyState == 4 && this.status == 200) {

        var responseObj = JSON.parse(this.responseText);
        var firstSoundCloudURL = responseObj.Songs[0];

        // The SoundCloud URLs for individual tracks have to be added as a HTTP GET parameter
        // to the URL 'https://w.soundcloud.com/player/'
        // Therefore, if we want to embed a track from the url 'https://soundcloud.com/youngmoneywayne/bankaccount',
        // we need to create a link that looks something like...
        // https://w.soundcloud.com/player/?url=https://soundcloud.com/youngmoneywayne/bankaccount
        let soundcloudEmbedURLs = [];
        for (var i = 0; i < responseObj.Songs.length; i++) {
          soundcloudEmbedURLs[i] = "https://w.soundcloud.com/player/?url=" + responseObj.Songs[i];
        }

        let newState = {
          songs: soundcloudEmbedURLs
        };
        setStateAsync(newState);

      // Handle Error Received from Server
      } else if (this.readyState == 4 && this.status != 200) {
        alert("ERROR: PC LOAD LETTER");
      }
    };
    httpRequest.send();
  }

  render() {

    return(

      <div style={mainDivStyle}>

        <h1>Topher's Top 100</h1>

        <h2>Trending Tracks from Twitter...</h2>

        <PlayList songs={ this.state.songs } />

      </div>

    );
  }
}

export default App;
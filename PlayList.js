import React from 'react';

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

class PlayList extends React.Component {

  render() {

    if (this.props.songs.length < 1) {
      return(
        <div>
          <h4>Loading...</h4>
        </div>
      );
    } else {

      let tracks = [];

      for (var i = 0; i < this.props.songs.length; i++) {
        tracks[i] = <div>
                      <iframe 
                        src={this.props.songs[i]} 
                        style={soundcloudIframeStyle}>
                      </iframe>
                      <br />
                    </div>
      }

      return(

        <div>

          {tracks}

        </div>

      );
    }
  }
}

export default PlayList;
import style from './avatar.scss';
import React, {Component} from 'react';

class Avatar extends Component {
  constructor(props) {
    super(props);
    this.state = props.user;
  }

  render() {
    return (
      <img src={this.state.avatar_url} className={style.avatar} />
    )
  }
}

export default Avatar;

import style from './avatar.module.scss';
import React, {Component} from 'react';

class Avatar extends Component {
  constructor(props) {
    super(props);
    this.state = props.user;
  }

  render() {
    return (
      <img src={this.state.avatar_url} alt={this.state.nickname} className={`${style.avatar} ${style[[this.props.class]]}`} />
    )
  }
}

export default Avatar;

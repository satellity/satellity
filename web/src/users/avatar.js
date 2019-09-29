import style from './avatar.scss';
import React, {Component} from 'react';

class Avatar extends Component {
  constructor(props) {
    super(props);
    this.state = props.user;
  }

  render() {
    let klass = style.avatar;
    if (this.props.class == 'small') {
      klass = klass + ' ' + style.small;
    }
    return (
      <img src={this.state.avatar_url} alt={this.state.nickname} className={klass} />
    )
  }
}

export default Avatar;

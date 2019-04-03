import style from './show.scss';
import moment from 'moment';
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import API from '../api/index.js';
import TopicItem from '../topics/item.js';

class UserShow extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.state = {user: {user_id: props.match.params.id, nickname: '', biography: '', avatar_url: '', created_at: ''}, topics: []}
  }

  componentDidMount() {
    this.api.user.show(this.state.user.user_id).then((user) => {
      user.created_at = moment(user.created_at).format('l');
      this.setState({user: user});
    });
    this.api.user.topics(this.state.user.user_id).then((data) => {
      this.setState({topics: data});
    });
  }

  render() {
    return (
      <View state={this.state} />
    )
  }
}

const View = (props) => {
  const topics = props.state.topics.map((topic) => {
    return (
      <TopicItem topic={topic} key={topic.topic_id}/>
    )
  });

  return (
    <div className='container'>
      <aside className='section aside'>
        <div className={style.profile}>
          <img src={props.state.user.avatar_url} className={style.avatar} />
          <div className={style.name}>
            {props.state.user.nickname}
          </div>
          <div className={style.created}>
            Joined {props.state.user.created_at}
          </div>
        </div>
      </aside>
      <main className='section main'>
        <ul className={style.topics}>
          {topics}
        </ul>
      </main>
    </div>
  )
};

export default UserShow;

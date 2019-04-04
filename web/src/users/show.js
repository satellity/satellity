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
    this.state = {
      user: {},
      topics: [],
    }
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
    let state = this.state;
    const topics = state.topics.map((topic) => {
      return (
        <TopicItem topic={topic} key={topic.topic_id}/>
      )
    });

    const profile = (
      <div className={style.profile}>
        <img src={state.user.avatar_url} className={style.avatar} />
        <div className={style.name}>
          {state.user.nickname}
        </div>
        <div className={style.created}>
          Joined {state.user.created_at}
        </div>
      </div>
    );

    return (
      <div className='container'>
        <aside className='section aside'>
          {profile}
        </aside>
        <main className='section main'>
          <ul className={style.topics}>
            {topics}
          </ul>
        </main>
      </div>
    )
  }
}

export default UserShow;

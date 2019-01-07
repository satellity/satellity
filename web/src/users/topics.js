import style from './index.scss';
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import TimeAgo from 'react-timeago';
import API from '../api/index.js';
import ColorUtils from '../components/color.js';

class UserTopics extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    const user = this.api.user.me();
    this.color = new ColorUtils();
    this.state = {user: {user_id: '', nickname: '', biography: '', avatar_url: '', created_at: ''}, topics: []};
  }

  componentDidMount() {
    const id = this.props.match.params["id"];
    this.api.user.show(id, (resp) => {
      this.setState({user: resp.data});
    });
    this.api.user.topics(id, (resp) => {
      this.setState({topics: resp.data});
    });
  }

  render() {
    return <View state={this.state} color={this.color} />
  }
}

const View = (props) => {
  const topics = props.state.topics.map((topic) => {
    let comment = '';
    if (topic.comments_count > 0) {
      comment = (
        <div className={style.topic_comment}>
          <span className={style.comments_count} style={{backgroundColor: props.color.colour(topic.topic_id)}}> {topic.comments_count} </span>
        </div>
      )
    }
    return (
      <li className={style.topic} key={topic.topic_id}>
        <div className={style.topic_detail}>
          <h2 className={style.topic_title}>
            <Link to={`/topics/${topic.topic_id}`}>
              {topic.title}
            </Link>
          </h2>
          <span>{topic.category.name}</span> • <TimeAgo date={topic.created_at} />
        </div>
        {comment}
      </li>
    )
  });

  return (
    <div>
      <div className={style.user}>
        <img src={props.state.user.avatar_url} className={style.avatar} />
        <div className={style.info}>
          <h1>
            {props.state.user.nickname}
          </h1>
          <div>
            {props.state.user.biography}
          </div>
        </div>
      </div>

      <div className='container'>
        <main className='section main'>
          <ul className={style.topics}>
            {topics}
          </ul>
        </main>
        <aside className='section aside'>
        </aside>
      </div>
    </div>
  )
}

export default UserTopics;

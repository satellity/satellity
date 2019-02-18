import style from './show.scss';
import topicStyle from '../styles/topic_item.scss';
import moment from 'moment';
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import TimeAgo from 'react-timeago';
import API from '../api/index.js';
import ColorUtils from '../components/color.js';

class UserShow extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.color = new ColorUtils();
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
      <View state={this.state} color={this.color} />
    )
  }
}

const View = (props) => {
  const topics = props.state.topics.map((topic) => {
    let comment = '';
    if (topic.comments_count > 0) {
      comment = (
        <span className={topicStyle.count} style={{backgroundColor: props.color.colour(topic.topic_id)}}> {topic.comments_count} </span>
      )
    }
    return (
      <li className={topicStyle.topic} key={topic.topic_id}>
        <div className={topicStyle.detail + ' ' + topicStyle.no_avatar}>
          <h2 className={topicStyle.title}>
            <Link to={`/topics/${topic.topic_id}`}>
              {topic.title}
            </Link>
          </h2>
          <span>{topic.category.name}</span> • <TimeAgo date={topic.created_at} />
        </div>
        <div className={topicStyle.comment}>
          {comment}
        </div>
      </li>
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

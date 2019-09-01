import style from './item.scss';
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import TimeAgo from 'react-timeago';
import ColorUtils from '../components/color.js';
import Avatar from '../users/avatar.js';

class TopicItem extends Component {
  constructor(props) {
    super(props);

    this.state = {
      profile: !!props.profile,
      topic: props.topic
    };
    this.color = new ColorUtils();
  }

  render() {
    let topic = this.state.topic;
    let comment;
    if (topic.comments_count > 0) {
      comment = (
        <span className={style.count} style={{backgroundColor: this.color.colour(topic.topic_id)}}> {topic.comments_count} </span>
      )
    }
    return (
      <li className={style.topic} key={topic.topic_id}>
          {!this.state.profile && <Avatar user={topic.user} />}
        <div className={style.detail}>
          <Link to={`/topics/${topic.short_id}-${topic.title.replace(/\W+/mgsi, ' ').replace(/\s+/mgsi, '-').replace(/[^\w-]/mgsi, '')}`}>
            <h2 className={style.title}>
              {topic.title}
            </h2>
          </Link>
          <div>
            {
              !this.state.profile &&
              <span>
                <Link to={`/users/${topic.user.user_id}`}>{topic.user.nickname.slice(0,16)}</Link>
              </span>
            }
            <span className={`${style.sep} ${this.state.profile ? style.left : ''}`}>{i18n.t('topic.in')}</span>
            <Link to={{pathname: "/", search: `?c=${topic.category.name}`}}>{topic.category.alias}</Link>
            <span className={style.sep}>{i18n.t('topic.at')}</span>
            <TimeAgo date={topic.created_at} />
          </div>
        </div>
        <div className={style.comment}>
          {comment}
        </div>
      </li>
    )
  }
}

export default TopicItem;

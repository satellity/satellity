import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import TimeAgo from 'react-timeago';
import style from './item.scss';
import ColorUtils from '../components/color.js';

class TopicItem extends Component {
  constructor(props) {
    super(props);
    this.state = {topic: props.topic};
    this.color = new ColorUtils();
  }

  render() {
    let topic = this.state.topic;
    let comment = '';
    if (topic.comments_count > 0) {
      comment = (
        <span className={style.count} style={{backgroundColor: this.color.colour(topic.topic_id)}}> {topic.comments_count} </span>
      )
    }
    return (
      <li className={style.topic} key={topic.topic_id}>
        <img src={topic.user.avatar_url} className={style.avatar} />
        <div className={style.detail}>
          <h2 className={style.title}>
            <Link to={`/topics/${topic.short_id}-${topic.title.replace(/\s+/mgsi, '-').replace(/[^\w-]/mgsi, '')}`}>{topic.title}</Link>
          </h2>
          <div>
            <span>{topic.user.nickname.slice(0,16)}</span>
            <span className={style.sep}>{i18n.t('topic.in')}</span><span>{topic.category.name}</span>
            <span className={style.sep}>{i18n.t('topic.at')}</span><TimeAgo date={topic.created_at} />
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

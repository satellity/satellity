import style from './item.module.scss';
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import TimeAgo from 'react-timeago';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import ColorUtils from '../components/color.js';
import Avatar from '../users/avatar.js';

class TopicItem extends Component {
  constructor(props) {
    super(props);

    this.state = {
      profile: !!props.profile,
      user: props.user,
      topic: props.topic
    };
    this.color = new ColorUtils();
  }

  render() {
    const i18n = window.i18n;
    let state = this.state;
    let topic = state.topic, comments;
    if (topic.comments_count > 0) {
      comments = (
        <span className={style.count}> {topic.comments_count} </span>
      )
    }
    return (
      <li className={style.topic} key={topic.topic_id}>
          {!this.state.profile && <Avatar user={topic.user} />}
        <div className={style.detail}>
          {
            topic.topic_type === 'POST' &&
            <Link to={`/topics/${topic.short_id}-${topic.title.replace(/\W+/mgsi, ' ').replace(/\s+/mgsi, '-').replace(/[^\w-]/mgsi, '')}`}>
              <h2 className={style.title}>
                  {topic.title}
                </h2>
            </Link>
          }
          {
            topic.topic_type === 'LINK' &&
            <a href={topic.body} target='_blank' rel='noopener noreferrer'>
              <h2 className={style.title}>
                  {topic.title} <FontAwesomeIcon icon={['fa', 'external-link-alt']} className={style.external} />
              </h2>
            </a>
          }
          <div>
            <Link className={style.node} to={`/categories/${topic.category.name}`} style={{color: this.color.colour(topic.category.name, 1), backgroundColor: this.color.colour(topic.category.name, 0.3)}}>{topic.category.alias}</Link>
            {
              !this.state.profile &&
                <span className={style.fullname}>
                <Link to={`/users/${topic.user.user_id}`}>{topic.user.nickname.slice(0,16)}</Link>
              </span>
            }
            <span className={style.sep}>{i18n.t('topic.at')}</span>
            <TimeAgo date={topic.created_at} />
            {topic.topic_type === 'LINK' && <Link to={`/topics/${topic.short_id}-${topic.title.replace(/\W+/mgsi, ' ').replace(/\s+/mgsi, '-').replace(/[^\w-]/mgsi, '')}`} className={style.comments}>{i18n.t('topic.comments')}</Link>}
          </div>
        </div>
        <div className={style.comment}>
          {comments}
        </div>
      </li>
    )
  }
}

export default TopicItem;

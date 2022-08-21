import React from 'react';
import {Link} from 'react-router-dom';
import TimeAgo from 'react-timeago';
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import ColorUtils from 'components/color.js';
import Avatar from 'users/avatar.js';
import PropTypes from 'prop-types';
import {seoTitle} from 'utils';

import style from './item.module.scss';

const TopicItem = (props) => {
  const {profile, topic} = props;

  const color = new ColorUtils();

  const i18n = window.i18n;
  let comments;
  if (topic.comments_count > 0) {
    comments = (
      <span className={style.count}> {topic.comments_count} </span>
    );
  }
  return (
    <li className={style.topic} key={topic.topic_id}>
      {!profile && <Avatar user={topic.user} />}
      <div className={style.detail}>
        {
          topic.topic_type === 'POST' &&
            <Link to={`/topics/${seoTitle(topic.title, topic.topic_id)}`}>
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
          <Link className={style.node} to={`/categories/${topic.category.name}`}
            style={{color: color.display(topic.category.name, 1), backgroundColor: color.display(topic.category.name, 0.3)}}>
            {topic.category.alias}
          </Link>
          {
            !profile &&
              <span className={style.fullname}>
                <Link to={`/users/${topic.user.user_id}`}>{topic.user.nickname.slice(0, 16)}</Link>
              </span>
          }
          <span className={style.sep}>{i18n.t('topic.at')}</span>
          <TimeAgo date={topic.created_at} />
          {
            topic.topic_type === 'LINK' &&
              <Link to={`/topics/${seoTitle(topic.title, topic.topic_id)}}`}
                className={style.comments}>
                {i18n.t('topic.comments')}
              </Link>
          }
        </div>
      </div>
      <div className={style.comment}>
        {comments}
      </div>
    </li>
  );
};

TopicItem.propTypes = {
  profile: PropTypes.bool,
  user: PropTypes.object,
  topic: PropTypes.object,
};

export default TopicItem;

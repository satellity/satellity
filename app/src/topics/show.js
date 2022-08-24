import style from './show.module.scss';
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import React, {useState, useEffect} from 'react';
import {Link, useParams} from 'react-router-dom';
import TimeAgo from 'react-timeago';
import showdown from 'showdown';
import showdownHighlight from 'showdown-highlight';
import {Helmet} from 'react-helmet';
import API from 'api/index.js';
import Config from 'components/config.js';
import SiteWidget from 'home/widget.js';
import CommentList from 'comments/index.js';
import Loading from 'components/loading.js';
import {titleToId} from 'utils';

const Show = () => {
  const converter = new showdown.Converter({extensions: ['header-anchors', showdownHighlight]});
  const {id} = useParams();
  const [loading, setLoading] = useState(true);
  const [i18n] = useState(window.i18n);
  const [api] = useState(new API());
  const [topic, setTopic] = useState({});

  useEffect(() => {
    setLoading(true);
    console.log(id, titleToId(id));
    api.topic.show(titleToId(id)).then((resp) => {
      if (resp.error) {
        return;
      }
      const data = resp.data;
      data.short_body = data.body.substring(0, 128);
      data.html_body = converter.makeHtml(data.body);
      setTopic(data);
      setLoading(false);
    });
  }, []);

  const handleClick = (e, action) => {
  };

  const loadingView = (
    <div className={style.loading}>
      <Loading class='medium' />
    </div>
  );

  let seoView;
  if (!loading) {
    seoView = (
      <Helmet>
        <title>{`${topic.title} - ${topic.user.nickname} - ${Config.Name}`}</title>
        <meta name='description' content={topic.short_body} />
        <link rel="canonical"
          href={`${Config.Host}/topics/${topic.short_id}-${topic.title.replace(/\W+/mgsi, ' ').replace(/\s+/mgsi, '-').replace(/[^\w-]/mgsi, '')}`} />
      </Helmet>
    );
  }

  let action;
  if (false) {
    action = (
      <Link to={`/topics/${topic.topic_id}/edit`} className={style.edit}>
        <FontAwesomeIcon icon={['far', 'edit']} />
      </Link>
    );
  }

  let like = {};
  if (topic.is_liked_by) {
    like = {color: 'rgb(218, 40, 16)'};
  }

  let bookmark = {};
  if (topic.is_bookmarked_by) {
    bookmark = {color: 'rgb(218, 40, 16)'};
  }

  let topicView;
  if (!loading) {
    topicView = (
      <div className={style.content}>
        <header className={style.header}>
          <div className={style.heading}>
            <h1>
              {topic.title}
              {action}
            </h1>
            <div className={style.info}>
              <Link to={`/users/${topic.user.user_id}`}>
                {topic.user.nickname}
              </Link>
              <span className={style.sep}>{i18n.t('topic.in')}</span>
              <Link to={{pathname: '/', search: `?c=${topic.category.name}`}}>{topic.category.alias}</Link>
              <span className={style.sep}>{i18n.t('topic.at')}</span>
              <TimeAgo date={topic.created_at} />
              <span className={style.views}>{topic.views_count} views</span>
            </div>
          </div>
          <img src={topic.user.avatar_url} className={style.avatar} alt={topic.user.nickname} />
        </header>
        <div>
          {topic.body !== '' && <article className={`md ${style.body}`} dangerouslySetInnerHTML={{__html: topic.html_body}} />}
        </div>
        <div className={style.actions}>
          {
            topic.is_owner &&
              <Link to={`/topics/${topic.topic_id}/edit`} className={style.action}>
                <FontAwesomeIcon icon={['far', 'edit']} />
              </Link>
          }
          <span className={style.item}>
            {
              topic.actioning !== 'like' &&
                <span className={`${style.action} ${topic.is_liked_by}`} onClick={(e) => handleClick(e, 'like')}>
                  {topic.likes_count > 0 && <span>{topic.likes_count}</span>}
                  <FontAwesomeIcon icon={['far', 'heart']} style={like}/>
                </span>
            }
            {topic.actioning === 'like' && <Loading class='small' />}
          </span>
          <span className={style.item}>
            {
              topic.actioning !== 'bookmark' &&
                <span className={`${style.action} ${topic.is_bookmarked_by}`} onClick={(e) => handleClick(e, 'bookmark')}>
                  {topic.bookmarks_count > 0 && <span>{topic.bookmarks_count}</span>}
                  <FontAwesomeIcon icon={['far', 'bookmark']} style={bookmark}/>
                </span>
            }
            {topic.actioning === 'bookmark' && <Loading class='small' />}
          </span>
        </div>
      </div>
    );
  }

  return (
    <div className='container'>
      {seoView}
      <main className='column main'>
        {loading && loadingView}
        {!loading && topicView}
        {!loading && <CommentList topicId={topic.topic_id} commentsCount={topic.comments_count} />}
      </main>
      <aside className='column aside'>
        <SiteWidget />
      </aside>
    </div>
  );
};

export default Show;

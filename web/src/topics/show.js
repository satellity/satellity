import style from './show.scss';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import TimeAgo from 'react-timeago';
import showdown from 'showdown';
import { Helmet } from 'react-helmet';
import API from '../api/index.js';
import Config from '../components/config.js';
import SiteWidget from '../home/widget.js';
import CommentList from '../comments/index.js';
import Loading from '../widgets/loading.js';

class Show extends Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: true,
      topic_id: props.match.params.id,
      user: {},
      category: {},
    };

    this.api = new API();
    this.converter = new showdown.Converter();
    this.handleClick = this.handleClick.bind(this);
  }

  componentDidMount() {
    const user = this.api.user.local();
    this.api.topic.show(this.props.match.params.id).then((resp) => {
      if (resp.error) {
        return
      }
      let data = resp.data;
      data.loading = false;
      data.is_owner = data.user.user_id === user.user_id;
      data.short_body = data.body.substring(0, 128);
      data.html_body = this.converter.makeHtml(data.body);
      this.setState(data);
    });
  }

  handleClick(e, action) {
    if (action === 'like' && this.state.is_liked_by) {
      action = 'unlike';
    }
    if (action === 'bookmark' && this.state.is_bookmarked_by) {
      action = 'unsave';
    }
    this.api.topic.action(action, this.state.topic_id).then((resp) => {
      if (resp.error) {
        return
      }
      this.setState(resp.data);
    });
  }

  render() {
    let state = this.state;
    const loadingView = (
      <div className={style.loading}>
        <Loading class='medium' />
      </div>
    )

    const seoView = (
      <Helmet>
        <title>{state.title} - {state.user.nickname} - {Config.Name}</title>
        <meta name='description' content={state.short_body} />
      </Helmet>
    )

    let action;
    if (state.is_owner) {
      action = (
        <Link to={`/topics/${state.topic_id}/edit`} className={style.edit}>
          <FontAwesomeIcon icon={['far', 'edit']} />
        </Link>
      )
    }

    let like = {};
    if (state.is_liked_by) {
    like = {color: 'rgb(218, 40, 16)'};
    }

    let bookmark = {};
    if (state.is_bookmarked_by) {
      bookmark = {color: 'rgb(218, 40, 16)'};
    }

    const topicView = (
      <div className={style.content}>
        <header className={style.header}>
          <div className={style.heading}>
            <h1>
              {state.title}
              {action}
            </h1>
            <div className={style.info}>
              {state.user.nickname}
              <span className={style.sep}>{i18n.t('topic.in')}</span>
              <Link to={{pathname: "/", search: `?c=${state.category.name}`}}>{state.category.alias}</Link>
              <span className={style.sep}>{i18n.t('topic.at')}</span>
              <TimeAgo date={state.created_at} />
            </div>
          </div>
          <img src={state.user.avatar_url} className={style.avatar} />
        </header>
        <div>
          {state.body !== '' && <article className={`md ${style.body}`} dangerouslySetInnerHTML={{__html: state.html_body}} />}
        </div>
        <div className={style.actions}>
          <span className={`${style.action} ${state.is_liked_by}`} onClick={(e) => this.handleClick(e, 'like')}>
            {state.likes_count > 0 && <span>{state.likes_count}</span>}
            <FontAwesomeIcon icon={['far', 'heart']} style={like}/>
          </span>
          <span className={`${style.action} ${state.is_bookmarked_by}`} onClick={(e) => this.handleClick(e, 'bookmark')}>
            {state.bookmarks_count > 0 && <span>{state.bookmarks_count}</span>}
            <FontAwesomeIcon icon={['far', 'bookmark']} style={bookmark}/>
          </span>
        </div>
      </div>
    )

    return (
      <div className='container'>
        {!state.loading && seoView}
        <main className='column main'>
          {state.loading && loadingView}
          {!state.loading && topicView}
          {!state.loading && <CommentList topicId={state.topic_id} commentsCount={state.comments_count} />}
        </main>
        <aside className='column aside'>
          <SiteWidget />
        </aside>
      </div>
    )
  }
}

export default Show;

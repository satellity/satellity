import style from './show.module.scss';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import TimeAgo from 'react-timeago';
import showdown from 'showdown';
import showdownHighlight from 'showdown-highlight';
import { Helmet } from 'react-helmet';
import API from '../api/index.js';
import Config from '../components/config.js';
import SiteWidget from '../home/widget.js';
import CommentList from '../comments/index.js';
import Loading from '../components/loading.js';

class Show extends Component {
  constructor(props) {
    super(props);
    this.state = {
      actioning: '',
      loading: true,
      topic_id: props.match.params.id,
      user: {},
      category: {},
    };

    this.api = new API();
    this.converter = new showdown.Converter({ extensions: ['header-anchors', showdownHighlight] });
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
    this.setState({actioning: action});
    if (action === 'like' && this.state.is_liked_by) {
      action = 'unlike';
    }
    if (action === 'bookmark' && this.state.is_bookmarked_by) {
      action = 'unsave';
    }
    this.api.topic.action(action, this.state.topic_id).then((resp) => {
      if (resp.error) {
        this.setState({actioning: ''});
        return
      }
      resp.data.actioning = '';
      this.setState(resp.data);
    });
  }

  render() {
    const i18n = window.i18n;
    let state = this.state;
    const loadingView = (
      <div className={style.loading}>
        <Loading class='medium' />
      </div>
    )

    let seoView;
    if (!state.loading) {
      seoView = (
        <Helmet>
          <title>{`${state.title} - ${state.user.nickname} - ${Config.Name}`}</title>
          <meta name='description' content={state.short_body} />
          <link rel="canonical" href={`${Config.Host}/topics/${state.short_id}-${state.title.replace(/\W+/mgsi, ' ').replace(/\s+/mgsi, '-').replace(/[^\w-]/mgsi, '')}`} />
        </Helmet>
      )
    }

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
              <Link to={`/users/${state.user.user_id}`}>
                {state.user.nickname}
              </Link>
              <span className={style.sep}>{i18n.t('topic.in')}</span>
              <Link to={{pathname: "/", search: `?c=${state.category.name}`}}>{state.category.alias}</Link>
              <span className={style.sep}>{i18n.t('topic.at')}</span>
              <TimeAgo date={state.created_at} />
              <span className={style.views}>{state.views_count} views</span>
            </div>
          </div>
          <img src={state.user.avatar_url} className={style.avatar} alt={state.user.nickname} />
        </header>
        <div>
          {state.body !== '' && <article className={`md ${style.body}`} dangerouslySetInnerHTML={{__html: state.html_body}} />}
        </div>
        <div className={style.actions}>
          {
            state.is_owner &&
            <Link to={`/topics/${state.topic_id}/edit`} className={style.action}>
              <FontAwesomeIcon icon={['far', 'edit']} />
            </Link>
          }
          <span className={style.item}>
            {
              state.actioning !== 'like' &&
              <span className={`${style.action} ${state.is_liked_by}`} onClick={(e) => this.handleClick(e, 'like')}>
                {state.likes_count > 0 && <span>{state.likes_count}</span>}
                <FontAwesomeIcon icon={['far', 'heart']} style={like}/>
              </span>
            }
            {state.actioning === 'like' && <Loading class='small' />}
          </span>
          <span className={style.item}>
            {
              state.actioning !== 'bookmark' &&
              <span className={`${style.action} ${state.is_bookmarked_by}`} onClick={(e) => this.handleClick(e, 'bookmark')}>
                {state.bookmarks_count > 0 && <span>{state.bookmarks_count}</span>}
                <FontAwesomeIcon icon={['far', 'bookmark']} style={bookmark}/>
              </span>
            }
            {state.actioning === 'bookmark' && <Loading class='small' />}
          </span>
        </div>
      </div>
    )

    return (
      <div className='container'>
        {seoView}
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

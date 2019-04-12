import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import TimeAgo from 'react-timeago';
import showdown from 'showdown';
import { Helmet } from 'react-helmet';
import API from '../api/index.js';
import style from './style.scss';
import SiteWidget from '../home/widget.js';
import CommentList from '../comments/index.js';
import LoadingView from '../loading/loading.js';

class TopicShow extends Component {
  constructor(props) {
    super(props);

    this.api = new API();
    this.converter = new showdown.Converter();
    this.state = {
      loading: true,
      topic_id: props.match.params.id,
      user: {},
      category: {},
    };
  }

  componentDidMount() {
    const user = this.api.user.readMe();
    this.api.topic.show(this.props.match.params.id).then((data) => {
      data.loading = false;
      data.is_author = data.user.user_id === user.user_id;
      data.short_body = data.body.substring(0, 128);
      data.body = this.converter.makeHtml(data.body);
      this.setState(data);
    });
  }

  render() {
    let state = this.state;
    const loadingView = (
      <div className={style.loading}>
        <LoadingView style='md-ring' />
      </div>
    )

    const seoView = (
      <Helmet>
        <title>{state.title} - {state.user.nickname} - GoDiscourse</title>
        <meta name='description' content={state.short_body} />
      </Helmet>
    )

    let editAction;
    if (state.is_author) {
      editAction = (
        <Link to={`/topics/${state.topic_id}/edit`} className={style.edit}>
          <FontAwesomeIcon icon={['far', 'edit']} />
        </Link>
      )
    }

    const topicView = (
      <div className={style.content}>
        <header className={style.header}>
          <div className={style.heading}>
            <h1>
              {state.title}
              {editAction}
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
          {state.body !== '' && <article className={`md ${style.body}`} dangerouslySetInnerHTML={{__html: state.body}} />}
        </div>
      </div>
    )

    return (
      <div className='container'>
        {!state.loading && seoView}
        <main className='section main'>
          {state.loading && loadingView}
          {!state.loading && topicView}
          {!state.loading && <CommentList topicId={state.topic_id} commentsCount={state.comments_count} />}
        </main>
        <aside className='section aside'>
          <SiteWidget />
        </aside>
      </div>
    )
  }
}

export default TopicShow;

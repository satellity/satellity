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
      topic_id: props.match.params.id, title: '', short_body: '', body: '', comments_count: 0, created_at: '', user_id: '', is_author: false, loading: true,
      user: {user_id: '', nickname: ''},
      category: {category_id: '', name: ''}
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
    return (
      <View state={this.state} />
    )
  }
}

const View = ({state}) => {
  let editAction = '';
  if (state.is_author) {
    editAction = (
      <Link to={`/topics/${state.topic_id}/edit`} className={style.edit}>
        <FontAwesomeIcon icon={['far', 'edit']} />
      </Link>
    )
  }
  const loadingView = (
    <div className={style.loading}>
      <LoadingView style='md-ring' />
    </div>
  )
  return (
    <div className='container'>
      <Helmet>
        <title>{state.title} - {state.user.nickname} - GoDiscourse</title>
        <meta name='description' content={state.short_body} />
      </Helmet>
      <main className='section main'>
        {state.loading && loadingView}
        <div className={style.content}>
          <header className={style.header}>
            <img src={state.user.avatar_url} className={style.avatar} />
            <h1 className={style.title}>
              {state.title}
              {editAction}
            </h1>
            <div className={style.info}>
              {state.category.name} • {state.user.nickname} • <TimeAgo date={state.created_at} />
            </div>
          </header>
          {state.body !== '' && <article className={`md ${style.body}`} dangerouslySetInnerHTML={{__html: state.body}}>
          </article>}
        </div>
        {state.title !== "" && <CommentList topicId={state.topic_id} commentsCount={state.comments_count} />}
      </main>
      <aside className='section aside'>
        <SiteWidget />
      </aside>
    </div>
  )
}

export default TopicShow;

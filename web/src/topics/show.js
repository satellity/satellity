import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import TimeAgo from 'react-timeago';
import showdown from 'showdown';
import API from '../api/index.js';
import style from './style.css';
import SiteWidget from '../components/site-widget.js';
import CommentIndex from '../comments/index.js';
import CommentNew from '../comments/new.js';

class TopicShow extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.converter = new showdown.Converter();
    this.state = {
      topic_id: props.match.params.id, title: '', body: '', comments_count: 0, created_at: '', user_id: '', is_author: false,
      user: {user_id: '', nickname: ''},
      category: {category_id: '', name: ''}
    };
    const classes = document.body.classList.values();
    document.body.classList.remove(...classes);
    document.body.classList.add('topic', 'layout');
  }

  componentDidMount() {
    const user = this.api.user.me();
    this.api.topic.show(this.props.match.params.id, (resp) => {
      resp.data.is_author = resp.data.user.user_id === user.user_id;
      resp.data.body = this.converter.makeHtml(resp.data.body);
      this.setState(resp.data);
    });
  }

  render() {
    return (
      <View state={this.state} />
    )
  }
}

const View = ({state}) => {
  var editAction = '';
  if (state.is_author) {
    editAction = (
      <Link to={`/topics/${state.topic_id}/edit`} className={style.edit}>
        <FontAwesomeIcon icon={['far', 'edit']} />
      </Link>
    )
  }
  return (
    <div className='container'>
      <main className='section main'>
        <header className={style.header}>
          <h1 className={style.title}>
            <img src={state.user.avatar_url} className={style.avatar} />
            {state.title}
            {editAction}
          </h1>
          <div className={style.detail}>
            {state.category.name} • {state.user.nickname} • <TimeAgo date={state.created_at} />
          </div>
        </header>
        <article className={`md ${style.body}`} dangerouslySetInnerHTML={{__html: state.body}}>
        </article>
        <CommentIndex topicId={state.topic_id} />
        <CommentNew topicId={state.topic_id} />
      </main>
      <aside className='section aside'>
        <SiteWidget />
      </aside>
    </div>
  )
}

export default TopicShow;

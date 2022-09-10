import React, {useEffect, useState} from 'react';
import {Navigate, useParams} from 'react-router-dom';
import CodeMirror from '@uiw/react-codemirror';
import {markdown, markdownLanguage} from '@codemirror/lang-markdown';
import {languages} from '@codemirror/language-data';
import {validate} from 'uuid';
import {useCategory} from 'services';
import API from 'api/index.js';
import Loading from 'components/loading.js';
import Button from 'components/button.js';

import style from './new.module.scss';

const Form = (props) => {
  const [api] = useState(new API());
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [topic] = useState({});
  const [topicId, setTopicId] = useState('draft');
  const [categoryId, setCategoryId] = useState('');
  const [title, setTitle] = useState('');
  const [body, setBody] = useState('');
  // const [draft] = useState(false);

  const {id} = useParams();

  useEffect(() => {
    if (id) {
      setTopicId(id);
    }
    api.topic.show(topicId).then((resp) => {
      if (resp.error) {
        return;
      }
      if (resp.data) {
        setTopicId(resp.data.topic_id);
        setCategoryId(resp.data.category_id);
        setTitle(resp.data.title);
      }
      setLoading(false);
    });
  }, [id]);

  const handleChange = (e) => {
    setTitle(e.target.value);
  };

  const handleBodyChange = (value) => {
    setBody(value);
  };

  const handleCategoryClick = (e, value) => {
    setCategoryId(value);
  };

  const submitForm = (e) => {
    e.preventDefault();
    if (submitting) {
      return;
    }
    setSubmitting(true);
    let request = api.topic.create(topic);
    if (validate(topic.topic_id)) {
      request = api.topic.update(topic.topic_id, topic);
    }
    request.then((resp) => {
      if (resp.error) {
        return;
      }
      setSubmitting(false);
    });
  };


  const {isLoading, data} = useCategory();

  useEffect(() => {
    if (data.length > 0) {
      setCategoryId(data[0].category_id);
    };
  }, [data]);

  let categories = [];
  if (!isLoading) {
    categories = data.map((c) => {
      return (
        <span key={c.category_id} className={`${style.category} ${c.category_id === categoryId ? style.active : ''}`}
          onClick={(e) => handleCategoryClick(e, c.category_id)}>
          {c.alias}
        </span>
      );
    });
  }

  if (loading) {
    return (
      <div className={style.loading}>
        <Loading class='medium'/>
      </div>
    );
  }

  let titleView = <h1>{i18n.t('topic.title.new')}</h1>;
  if (validate(topicId)) {
    titleView = <h1>{i18n.t('topic.title.edit', {name: title})}</h1>;
  }

  const form = (
    <form onSubmit={submitForm}>
      <div className={style.categories}>
        {categories}
      </div>
      <div>
        <input type='text' name='title' pattern='.{3,}' required value={title} autoComplete='off' placeholder={i18n.t('topic.placeholder.title')}
          onChange={handleChange} />
      </div>
      <div>
        <CodeMirror
          value={body}
          height={body.split('\n').length > 16 ? 'auto': '300px'}
          extensions={[markdown({base: markdownLanguage, codeLanguages: languages})]}
          onChange={handleBodyChange}
        />
      </div>
      <div className={style.upload}>
        <a href='https://imgur.com/upload' target='_blank' rel='noopener noreferrer'>Choose imgur upload image first.</a>
      </div>
      <div className={style.submit}>
        <Button type='submit' classes='submit' disabled={submitting} text={i18n.t('general.submit')} />
      </div>
    </form>
  );

  return (
    <>
      {titleView}
      {form}
    </>
  );
};

const New = () => {
  const api = new API();
  if (!api.me.value()) {
    return (
      <Navigate to="/" replace />
    );
  }

  return (
    <div className='container'>
      <main className='column main'>
        <div className={style.form}>
          <Form />
        </div>
      </main>
      <aside className='column aside'>
        <div className={style.title}>Rules</div>
        <ul className={style.rules} dangerouslySetInnerHTML={{__html: i18n.t('topic.rules')}}></ul>
      </aside>
    </div>
  );
};

// const TOOLBAR = [
//   {icon: 'heading', action: 'heading', identity: ''},
//   {icon: 'bold', action: 'bold', identity: '**'},
//   {icon: 'italic', action: 'italic', identity: '*'},
//   {icon: 'strikethrough', action: 'strikethrough', identity: '~~'},
//   {icon: 'quote-left', action: 'quote', identity: '> '},
//   {icon: 'list-ol', action: 'ol', identity: '1. '},
//   {icon: 'list-ul', action: 'ul', identity: '* '},
// ];

export default New;

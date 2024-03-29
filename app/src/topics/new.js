import React, {useEffect, useState} from 'react';
import {Navigate, useParams, useNavigate} from 'react-router-dom';
import CodeMirror from '@uiw/react-codemirror';
import {markdown, markdownLanguage} from '@codemirror/lang-markdown';
import {languages} from '@codemirror/language-data';
import {validate} from 'uuid';
import {useCategory} from 'services';
import API from 'api/index.js';
import Loading from 'components/loading.js';
import Button from 'components/button.js';
import {seoTitle} from 'utils';

import style from './new.module.scss';

const Nodes = ({categoryId, setCategoryId}) => {
  const {isLoading, data} = useCategory();

  const handleCategoryClick = (e, value) => {
    setCategoryId(value);
  };

  if (isLoading) {
    return;
  }

  useEffect(() => {
    if (data.length > 0 && !categoryId) {
      setCategoryId(data[0].category_id);
    };
  }, [data]);

  const categories = data.map((c) => {
    return (
      <span key={c.category_id} className={`${style.category} ${c.category_id === categoryId ? style.active : ''}`}
        onClick={(e) => handleCategoryClick(e, c.category_id)}>
        {c.alias}
      </span>
    );
  });
  return categories;
};

const Form = (props) => {
  const navigate = useNavigate();
  const {id} = useParams();

  const [api] = useState(new API());
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [topicId, setTopicId] = useState('draft');
  const [categoryId, setCategoryId] = useState('');
  const [title, setTitle] = useState('');
  const [body, setBody] = useState('');
  const [draft, setDraft] = useState(true);

  useEffect(() => {
    api.topic.show(id || topicId).then((resp) => {
      if (resp.error) {
        return;
      }
      const topic = resp.data;
      if (topic) {
        setTopicId(topic.topic_id);
        setCategoryId(topic.category_id);
        setTitle(topic.title);
        setBody(topic.body);
        setDraft(topic.draft);
      }
      setLoading(false);
    });
  }, [id]);

  useEffect(() => {
    const handler = setTimeout(() => {
      submitForm(undefined, true);
    }, 10000);
    return () => {
      clearTimeout(handler);
    };
  }, [categoryId, title, body]);

  const handleChange = (e) => {
    setTitle(e.target.value);
  };

  const handleBodyChange = (value) => {
    setBody(value);
  };

  const submitForm = (e, autoSave) => {
    if (e) {
      e.preventDefault();
    }
    if (submitting) {
      return;
    }
    if (!validate(categoryId)) {
      return;
    }
    if (title.length <= 3) {
      return;
    }
    setSubmitting(true);
    let dra = draft;
    if (!autoSave) {
      dra = false;
    }
    const params = {title, body, category_id: categoryId, topic_type: 'POST', draft: dra};
    let request;
    if (validate(topicId)) {
      request = api.topic.update(topicId, params);
    } else {
      request = api.topic.create(params);
    }
    request.then((resp) => {
      if (resp.error) {
        return;
      }
      const topic = resp.data;
      setTopicId(topic.topic_id);
      if (!autoSave) {
        navigate(`/topics/${seoTitle(topic.title, topic.topic_id)}`);
      }
    }).finally(() => {
      setSubmitting(false);
    });
  };

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
    <form onSubmit={(e) => submitForm(e, false)}>
      <div className={style.categories}>
        <Nodes categoryId={categoryId} setCategoryId={setCategoryId}/>
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

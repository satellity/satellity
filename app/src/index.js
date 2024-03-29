import './index.scss';
import React from 'react';
import {createRoot} from 'react-dom/client';
import * as serviceWorker from './serviceWorker';
import {BrowserRouter, Route, Routes} from 'react-router-dom';
import showdown from 'showdown';
import {QueryClient, QueryClientProvider} from '@tanstack/react-query';
import API from './api/index.js';
import Locale from './locale/index.js';
import MainLayout from './layouts/main.js';
import NoMatch from './layouts/sink.js';
import AdminRoute from 'admin/admin.js';
import Admin from 'admin/pages.js';
import Oauth from './users/oauth.js';
import User from 'users/view.js';
import Topic from 'topics/view.js';
import Gist from 'gists/view.js';
import {library} from '@fortawesome/fontawesome-svg-core';
import {faBookmark, faComment, faComments, faEdit, faEye, faEyeSlash, faTrashAlt, faHeart} from '@fortawesome/free-regular-svg-icons';
import {
  faChalkboard, faEllipsisV, faHome,
  faPlus, faUsersCog, faExternalLinkAlt,
  faHeading, faBold, faItalic, faStrikethrough,
  faQuoteLeft, faListOl, faListUl, faTags,
} from '@fortawesome/free-solid-svg-icons';
import {faMarkdown} from '@fortawesome/free-brands-svg-icons';
library.add(
    faBookmark, faComment, faComments,
    faEdit, faEye, faEyeSlash, faTrashAlt,
    faHeart, faExternalLinkAlt,
    faChalkboard, faEllipsisV, faHome,
    faPlus, faUsersCog, faMarkdown,
    faHeading, faBold, faItalic, faStrikethrough, faQuoteLeft, faListOl, faListUl,
    faTags,
);

showdown.setOption('customizedHeaderId', true);
showdown.setOption('simplifiedAutoLink', true);
showdown.setOption('openLinksInNewWindow', true);
showdown.setOption('strikethrough', true);
showdown.setOption('simpleLineBreaks', true);
showdown.setOption('parseImgDimensions', true);
showdown.setOption('omitExtraWLInCodeBlocks', true);
showdown.setOption('tables', true);
// TODO
showdown.extension('header-anchors', function() {
  // eslint-disable-next-line max-len
  const ancTpl = '$1$4<a id="user-content-$3" class="anchor" href="#$3" aria-hidden="true"><svg aria-hidden="true" class="octicon octicon-link" height="16" version="1.1" viewBox="0 0 16 16" width="16"><path fill-rule="evenodd" d="M4 9h1v1H4c-1.5 0-3-1.69-3-3.5S2.55 3 4 3h4c1.45 0 3 1.69 3 3.5 0 1.41-.91 2.72-2 3.25V8.59c.58-.45 1-1.27 1-2.09C10 5.22 8.98 4 8 4H4c-.98 0-2 1.22-2 2.5S3 9 4 9zm9-3h-1v1h1c1 0 2 1.22 2 2.5S13.98 12 13 12H9c-.98 0-2-1.22-2-2.5 0-.83.42-1.64 1-2.09V6.25c-1.09.53-2 1.84-2 3.25C6 11.31 7.55 13 9 13h4c1.45 0 3-1.69 3-3.5S14.5 6 13 6z"></path></svg></a>$5';

  return [{
    type: 'html',
    regex: /(<h([1-3]) id="([^"]+?)">)(.*)(<\/h\2>)/g,
    replace: ancTpl,
  }];
});

let language = navigator.language;
const locale = new URLSearchParams(window.location.search).get('locale');
if (!!locale) {
  language = locale.split('-')[0];
}
window.i18n = new Locale(language);
window.api = new API();

const queryClient = new QueryClient();

const container = document.getElementById('root');
const root = createRoot(container); // createRoot(container!) if you use TypeScript

root.render(
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          <Route path='/oauth/:provider/callback' element={<Oauth />} />
          <Route path='/admin' element={<AdminRoute />}>
            <Route index element={<Admin.Index />} />
            <Route exact path={`users`} element={<Admin.Users />} />
            <Route exact path={`topics`} element={<Admin.Topics />} />
            <Route exact path={`comments`} element={<Admin.Comments />} />
            <Route exact path={`categories`} element={<Admin.Category.Index />} />
            <Route exact path={`categories/new`} element={<Admin.Category.New />} />
            <Route path={`categories/:id/edit`} element={<Admin.Category.Edit />} />
            <Route exact path={`gists`} element={<Admin.Gists />} />
            <Route exact path={`sources`} element={<Admin.Sources />} />
          </Route>
          <Route path='/' element={<MainLayout />}>
            <Route index element={<Gist.Index />} />
            <Route path="categories/:id" element={<Topic.Index />} />
            <Route path='user/edit' element={<User.Edit />} />
            <Route path='users/:id' element={<User.Show />} />
            <Route path='topics/new' element={<Topic.New />} />
            <Route path='topics/:id/edit' element={<Topic.New />} />
            <Route path='topics/:id' element={<Topic.Show />} />
          </Route>
          <Route path='/404' element={<NoMatch />} />
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>,
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();

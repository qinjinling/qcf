import Root from './pages/Root.js'

m.route(document.body, '/', {
    '/':  Root,
    '/:page': Root
})
import Header from '../components/Header.js'
import SearchView from './SearchView.js'
import PartDetail from './PartDetail.js'

export default function Root() {
    
    function renderView(vnode) {
        if (vnode.attrs.page == 'detail') {
            return [m(Header), m(PartDetail)]
        }
        return [m(Header), m(SearchView)]
    }

    return {
        view: renderView
    }
}
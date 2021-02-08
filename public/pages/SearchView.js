export default function SearchView() {
    let pn = ''
    let disabled = false
    let clicked = false

    function PartNumberLink(row) {
        if (row.IsStock == 'No')
            return row.PartNumber;
        return m('a', {
            href: '#!/detail?pn=' + encodeURIComponent(row.PartNumber) + '&mpid=' + row.MasterPartID
        }, row.PartNumber)
    }

    function searchResult() {
        if (window.tables && window.tables.length > 0) {
            return m('ul', window.tables.map(function (row) {
                return m('li', {
                    style: row.IsStock == 'Yes' ? 'color:green;' : 'color:red;'
                }, PartNumberLink(row))
            }))
        } else {
            if (clicked) {
                return m('p', '没有查询到任何数据')
            }
        }
    }

    function search() {
        disabled = true
        m.request({
            method: "GET",
            url: "/search",
            params: {
                pn: pn
            },
        }).then(function (data) {
            disabled = false
            clicked = true
            window.tables = data.Tables
        })
    }

    function renderView() {
        return [
            m('input', {
                name: 'pn',
                type: 'text',
                style: 'width:300px',
                oninput: function () {
                    pn = this.value
                }
            }),

            m('button', {
                type: 'button',
                disabled: disabled,
                onclick: search
            }, 'Search'),

            searchResult()
        ]
    }

    return {
        view: renderView
    }
}
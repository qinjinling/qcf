export default function PartDetail() {
    let partNumber, mpid
    let detail, error

    function StripPnForUrlRoute(pn) {
        pn = pn.replace(/['/','\']/g, "-").replace(/[+,&,*,:,%,?]/g, "")
        return pn
    }

    function PartTable(data) {
        let showAnyManufacturer = false
        return m('table', [
            m('thead', m('tr', [
                m('th', '制造商'),
                m('th', '库存'),
                m('th', '发货日期'),
                m('th', { style: 'text-align:right;' }, '数量阶梯'),
                m('th', '美元')
            ])),
            m('tbody', data.map(function (row) {
                var borderStyle = 'border-bottom: 2px solid rgb(26, 188, 156);'
                var displyNone = 'display:none;'
                var rowStyle = row.BorderBottomHighligt ? borderStyle : ''
                if (showAnyManufacturer && row.IsAnyManufacturerInStock) {
                    rowStyle += displyNone;
                }
                return m('tr', { style: rowStyle }, [
                    m('td', row.IsAnyManufacturer ? [row.Manufacturer, m('small', m('button', {
                        onclick: function () {
                            showAnyManufacturer = !showAnyManufacturer
                        }
                    }, showAnyManufacturer ? '⇣查看厂家' : '⇡隐藏厂家'))] : row.Manufacturer),
                    m('td', row.AvailableQty),
                    m('td', row.ShipDate),
                    m('td', m('div', {
                        style: 'display:flex;flex-direction:column;text-align:right'
                    }, m.trust(row.Quantity))),
                    m('td', m('div', {
                        style: 'display:flex;flex-direction:column;'
                    }, m.trust(row.Pricing))),
                ])
            }))
        ])
    }

    function renderView() {
        if (error) {
            return [
                m('div', {
                    style: 'color: red'
                }, error),
                m('a', { href: '#!/' }, '返回')
            ]
        }
        if (!detail) {
            return m('h4', 'Loading...')
        }
        return [
            m('h1', { style: 'color:#3498db' }, [
                partNumber,
                m('small', { style: 'color: royalblue;font-size: 1rem;margin-left:5px;' }, detail.Category)
            ]),
            m('p', detail.Description),
            PartTable(detail.Details || []),
            m('footer', m('a', { href: '#!/' }, '返回'))
        ]
    }

    function init() {
        partNumber = m.route.param('pn')
        m.request({
            method: 'GET',
            url: '/detail',
            params: {
                pn: StripPnForUrlRoute(partNumber),
                mpid: m.route.param('mpid')
            },
        }).then(function (data) {
            detail = data
        }).catch(function (e) {
            error = e.message
        })
    }

    return {
        oninit: init,
        view: renderView
    }
}
(function () {
  var root = document.body
  var clicked = false
  var ds = []

  var StripPnForUrlRoute = function (pn) {
    pn = pn.replace(/['/','\']/g, "-").replace(/[+,&,*,:,%,?]/g, "");
    return pn;
  }

  var Header = {
    view: function () {
      return m('h1', 'questcomp.com Facker')
    }
  }

  // 查询组件
  var SearchBar = {
    pn: '',
    disabled: false,
    view: function () {
      var self = this
      var input = m('input', {
        name: 'pn',
        type: 'text',
        style: 'width:300px',
        oninput: function () {
          self.pn = this.value
        }
      })
      var button = m('button', {
        type: 'button',
        disabled: self.disabled,
        onclick: function () {
          self.disabled = true
          m.request({
            method: "POST",
            url: "/search?pn=" + self.pn,
            withCredentials: true,
          })
          .then(function (data) {
            self.disabled = false
            clicked = true
            ds = data.Tables
          })
        }
      }, 'Search')
      return [input, button]
    }
  }

  // 零件链接
  var PartNumberLink = {
    view: function (vnode) {
      var data = vnode.attrs
      if (data.IsStock == 'No')
        return data.PartNumber;
      return m('a', {
        href: '#!/detail/' + encodeURIComponent(data.PartNumber) +'/' + data.MasterPartID
      }, data.PartNumber)
    }
  }

  // 查询结果
  var SearchResult = {
    view: function (vnode) {
      var self = this
      var data = vnode.attrs.data
      if (data && data.length > 0) {
        return m('ul', data.map(function (row) {
          var link = m(PartNumberLink, row)
          return m('li', {
            style: row.IsStock == 'Yes' ? 'color:green;' : 'color:red;'
          }, link)
        }))
      } else {
        if (clicked) {
          return m('p', '没有查询到任何数据')
        }
      }
    }
  }

  var Home = {
    view: function () {
      return [
        m(Header),
        m(SearchBar),
        m(SearchResult, { data: ds })
      ]
    }
  }

  var PartTable = {
    view: function (vnode) {
      var tableHeading = m('tr', [
        m('th', 'Manufacturer'),
        m('th', 'Available Qty'),
        m('th', 'Ship Date'),
        m('th', {style: 'text-align:right;'}, 'Quantity'),
        m('th', 'Pricing(USD)')
      ])
      var tableRows = vnode.attrs.data.map(function (row) {
        return m('tr', {
          style: row.BorderBottomHighligt ? 'border-bottom: 2px solid rgb(26, 188, 156);' : ''
        }, [
          m('td', row.Manufacturer),
          m('td', row.AvailableQty),
          m('td', row.ShipDate),
          m('td', m('div', {
            style: 'display:flex;flex-direction:column;text-align:right'
          }, m.trust(row.Quantity))),
          m('td', m('div', {
            style: 'display:flex;flex-direction:column;'
          }, m.trust(row.Pricing))),
        ])
      })
      return m('table', {
        cellspacing: '0'
      }, [
        m('thead', tableHeading),
        m('tbody', tableRows)
      ])
    }
  }

  var PartDetail = {
    oninit: function(vnode) {
      var self = this
      m.request({
        method: 'POST',
        url: '/detail?pn=' + StripPnForUrlRoute(vnode.attrs.pn) + '&mpid=' + vnode.attrs.mpid,
        withCredentials: true,
      })
      .then(function (data) {
        self.detail = data
      })
    },
    view: function (vnode) {
      var self = this
      if (!self.detail) {
        return [m(Header), m('h4', 'Loading...')]
      }
      return [
        m(Header),
        m('h1', {style: 'color:#3498db'}, [
          vnode.attrs.pn,
          m('small', {style: 'color: royalblue;font-size: 1rem;margin-left:5px;'}, self.detail.Category)
        ]),
        m('p', self.detail.Description),
        m(PartTable, { data: self.detail.Details }),
        m('footer', m('a', {href: '#!/'}, '返回'))
      ]
    }
  }

  m.route(root, '/', {
    '/': Home,
    '/detail/:pn/:mpid': PartDetail
  })
}())
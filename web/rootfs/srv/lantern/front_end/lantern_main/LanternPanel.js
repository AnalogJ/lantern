// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// function removePanel(view){
//     if (!UI.inspectorView._tabbedPane.hasTab(view.viewId()))
//         return;
//
//     delete view[UI.ViewManager._Location.symbol];
//     UI.viewManager._views.delete(view.viewId());
//     UI.inspectorView._tabbedPane.closeTab(view.viewId());
//     // this._views.delete(view.viewId());
// }


LanternMain.LanternPanel = class extends UI.Panel {
  constructor() {
    super('node-connection');
    this.registerRequiredCSS('lantern_main/nodeConnectionsPanel.css');
    this.contentElement.classList.add('node-panel');

    const container = this.contentElement.createChild('div', 'node-panel-center');

    const image = container.createChild('img', 'node-panel-logo');
    image.src = '/logo.png';


    /** @type {!Adb.Config} */
    this._config;

    this.contentElement.tabIndex = 0;
    this.setDefaultFocusedElement(this.contentElement);


    this._lanternView = new LanternMain.LanternView();
    this._lanternView.show(container);
  }
};

/**
 * @implements {UI.ListWidget.Delegate<Adb.PortForwardingRule>}
 */
LanternMain.LanternView = class extends UI.VBox {

  constructor(callback) {
    super();
    this._callback = callback;
    this.element.classList.add('network-discovery-view');

    const networkDiscoveryFooter = this.element.createChild('div', 'network-discovery-footer');
    networkDiscoveryFooter.createChild('span').textContent =
        Common.UIString('Peer into your network traffic. ');
    const link = networkDiscoveryFooter.createChild('span', 'link');
    link.textContent = Common.UIString('Learn more');
    link.addEventListener('click', () => InspectorFrontendHost.openInNewTab('https://github.com/analogj/lantern'));

    const inspectButton = UI.createTextButton(
        Common.UIString('Inspect HTTPS Traffic'), this._downloadCACertificateButtonClicked.bind(this), 'add-network-target-button',
        true /* primary */);
    this.element.appendChild(inspectButton);


    const mobileConfigButton = UI.createTextButton(
        Common.UIString('IOS MobileConfig'), this._downloadMobileConfigButtonClicked.bind(this), 'add-network-target-button',
        true /* primary */);
    this.element.appendChild(mobileConfigButton);

    this.element.classList.add('node-frontend');
  }

  _downloadCACertificateButtonClicked() {
      InspectorFrontendHost.openInNewTab('/certs/ca.crt')
  }
  _downloadMobileConfigButtonClicked() {
      InspectorFrontendHost.openInNewTab('/certs/lantern.mobileconfig')
  }
};

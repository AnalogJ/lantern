// Copyright 2018 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
/**
 * @implements {Common.Runnable}
 */
LanternMain.LanternMain = class extends Common.Object {
  /**
   * @override
   */
  run() {
    // Host.userMetrics.actionTaken(Host.UserMetrics.Action.ConnectToNodeJSFromFrontend);
    // SDK.initMainConnection(() => {
    //   const target = SDK.targetManager.createTarget('main', Common.UIString('Main'), SDK.Target.Type.Browser, null);
    //   target.setInspectedURL('Node.js');
    // }, Components.TargetDetachedDialog.webSocketConnectionLost);
  }
};


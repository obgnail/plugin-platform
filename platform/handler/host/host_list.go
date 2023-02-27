package host

//
//var hostList sync.Map
//
//func NewHostAndStore(pbHostInfo *protocol.HostDescriptor) types.IHost {
//	info, ok := hostList.Load(pbHostInfo.GetHostID())
//	if !ok {
//		hostInfo := pbToHost(pbHostInfo)
//		info = hostInfo
//	}
//	hostList.Store(pbHostInfo.GetHostID(), info.(*host))
//	return info.(*host)
//}
//
//func NewTmpHost(pbHostInfo *protocol.HostDescriptor) types.IHost {
//	info, ok := hostList.Load(pbHostInfo.GetHostID())
//	if !ok {
//		hostInfo := pbToHost(pbHostInfo)
//		info = hostInfo
//	}
//	return info.(*host)
//}
//
//func pbToHost(pbHostInfo *protocol.HostDescriptor) *host {
//	h := &host{}
//
//	h.info.IsLocal = pbHostInfo.GetIsLocal()
//	h.info.ID = pbHostInfo.GetHostID()
//	h.info.Name = pbHostInfo.GetName()
//
//	h.info.Version = buildmessageutils.GetVersionString(pbHostInfo.GetHostVersion())
//	h.info.SubVersion = buildmessageutils.GetVersionString(pbHostInfo.GetHostSubVersion())
//	h.info.MinSystemVersion = buildmessageutils.GetVersionString(pbHostInfo.GetMinSystemVersion())
//
//	h.info.Language = pbHostInfo.GetLanguage()
//	h.info.LanguageVersion = buildmessageutils.GetVersionString(pbHostInfo.GetLanguageVersion())
//	h.status = types.HostStatusNormal
//
//	return h
//}
//
//func Delete(hostID string) {
//	hostList.Delete(hostID)
//}
//
//func All() []types.IHost {
//	list := make([]types.IHost, 0)
//	hostList.Range(func(k, v interface{}) bool {
//		h := v.(types.IHost)
//		list = append(list, h)
//		return true
//	})
//	return list
//}
//
//func AllReleaseHost() []types.IHost {
//	list := make([]types.IHost, 0)
//	hostList.Range(func(k, v interface{}) bool {
//		h := v.(types.IHost)
//		if !h.GetInfo().IsLocal {
//			list = append(list, h)
//		}
//		return true
//	})
//	return list
//}
//
//func SyncWaitGetHostByLanguageAndHostVersion(language, hostVersion string) (hostInfo types.IHost, err error) {
//	beforeTime := time.Now()
//	info := make(chan types.IHost, 1)
//	loop := true
//	if len(All()) == 0 {
//		return nil, errors.Trace(fmt.Errorf("host not connected"))
//	}
//	go func(info chan types.IHost, loop *bool) {
//		for *loop {
//			hostList.Range(func(k, v interface{}) bool {
//				// 匹配host的版本以及host的语言
//				host := v.(types.IHost)
//				pluginHostVersionArray := strings.Split(hostVersion, ".")
//				prefix := strings.Join([]string{pluginHostVersionArray[0], pluginHostVersionArray[1]}, ".")
//
//				hostVersion := host.GetInfo().Version
//				if strings.HasPrefix(hostVersion, prefix) &&
//					host.GetInfo().Language == language {
//					info <- host
//					return false
//				}
//				return true
//			})
//			// 避免程序把cpu占满
//			time.Sleep(time.Millisecond * 200)
//		}
//	}(info, &loop)
//	select {
//	// 等待一个心跳多一点的时间
//	case <-time.After(time.Millisecond*time.Duration(config.Config.HeartBeatIntervalMs) + 100):
//		return nil, errors.Trace(fmt.Errorf("host not connected"))
//	case hostInfo = <-info:
//		loop = false
//		log.Logger.Trace("SyncWaitGetHostByLanguageAndHostVersion :%d", time.Since(beforeTime).Milliseconds())
//		return hostInfo, nil
//	}
//}
//
//func SyncWaitGetHostByHostID(hostID string) (hostInfo types.IHost, err error) {
//	checkIntervalMs := 100
//	retryTimes := (3 * int(config.Config.HeartBeatIntervalMs)) / checkIntervalMs
//	for i := 0; i < retryTimes; i++ {
//		if val, ok := hostList.Load(hostID); ok {
//			return val.(types.IHost), nil
//		}
//		time.Sleep(time.Millisecond * time.Duration(checkIntervalMs))
//		log.Logger.Trace("SyncWaitGetHostByHostID :sleep %d ms", checkIntervalMs)
//	}
//	return nil, errors.Trace(fmt.Errorf("host not connected "))
//}

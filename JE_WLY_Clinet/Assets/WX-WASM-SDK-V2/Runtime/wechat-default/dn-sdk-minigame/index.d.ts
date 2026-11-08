/**
 * SDK 初始化参数
 */
interface SdkConfig {
    /**
     * 数据源ID，必填
     */
    user_action_set_id: number;
    /**
     * 加密key，必填
     */
    secret_key: string;
    /**
     * 微信小程序/小游戏APPID，wx开头，必填
     */
    appid: string;
    /**
     * 微信 openid，openid 和 unionid 至少填一个, 可以调用 setOpenId 设置
     */
    openid?: string;
    /**
     * 微信 unionid，openid 和 unionid 至少填一个, 可以调用 setUnionId 设置
     */
    unionid?: string;
    /**
     * 自定义用户ID，选填
     */
    user_unique_id?: string;
    /**
     * 是否自动采集事件
     */
    auto_track?: boolean;
    /**
     * 是否自动采集属性
     */
    auto_attr?: boolean;
    /**
     * 上报回调（包括成功、失败、重试的各种场景）
     */
    on_report_complete?: Function | undefined;
    /**
     * 明确上报失败的回调（sdk内部会直接丢弃的行为，无法再重试上报成功的部分）
     */
    on_report_fail?: Function | undefined;
}
/**
 * 小程序/小游戏信息
 */
interface GameInfo {
    /**
     * 广告点击id、曝光id
     */
    ad_trace_id?: string;
    /**
     * 页面路径,小游戏没有
     */
    page_url?: string;
    /**
     * 页面名称,小游戏没有
     */
    page_title?: string;
    /**
     * 组件名称
     */
    component_name?: string;
    /**
     * 场景来源
     */
    source_scene?: number;
    /**
     * 渠道号
     */
    pkg_channel_id?: string;
    /**
     * 框架信息,小游戏没有
     */
    framework?: string;
    /**
     * 渠道来源
     */
    channel: string;
    /**
     * 启动参数
     */
    launch_options: string;
}
/**
 * 设备信息
 */
interface DeviceInfo {
    /**
     * 设备品牌
     */
    device_brand: string;
    /**
     * 设备型号
     */
    device_model: string;
    /**
     * 微信App版本
     */
    wx_version: string;
    /**
     * 微信基础库版本
     */
    wx_lib_version: string;
    /**
     * 微信客户端平台
     */
    wx_platform: string;
    /**
     * 操作系统
     */
    os: string;
    /**
     * 操作系统版本
     */
    os_version: string;
    /**
     * 设备性能等级
     */
    benchmark_level: number;
    /**
     * 屏幕高度
     */
    screen_height: number;
    /**
     * 屏幕宽度
     */
    screen_width: number;
}
interface TrackBaseInfo extends SdkConfig, DeviceInfo {
    local_id: string;
    sdk_name: string;
    sdk_version: string;
    inner_param: {
        app_version: string;
        app_env_version: string;
    };
}
/**
 * 行为信息
 */
interface ActionInfo {
    action_type: string;
    action_param?: Record<string | number | symbol, unknown>;
    is_sdk_auto_track: boolean;
    session_id: string;
    log_id: number;
    action_id: string;
    action_time: number;
    revised_action_time: number;
    is_retry: boolean;
    retry_count: number;
    ad_trace_id?: string;
    page_url?: string;
    page_title?: string;
    component_name?: string;
    source_scene?: number;
    pkg_channel_id?: string;
    /**
     * 网络类型
     */
    network_type?: string;
    inner_param?: Record<string | number | symbol, unknown>;
    channel?: string;
}
interface LogParams {
    log_type: string;
    code?: number;
    message?: string;
    err_stack?: string;
    user_action_set_id?: number;
    appid?: string;
    local_id?: string;
    session_id?: string;
    sdk_name?: string;
    sdk_version?: string;
    page_url?: string;
    framework?: string;
    device_brand?: string;
    device_model?: string;
    wx_version?: string;
    wx_lib_version?: string;
    wx_platform?: string;
    os?: string;
    os_version?: string;
}
interface TrackResponse {
    code: number;
    message: string;
    actionId?: string;
}

interface QueueProps {
    /**
     * 数据源ID，用于本地存储队列key
     */
    userActionSetId: number;
    /**
     * 队列长度
     */
    maxLength: number;
}
interface QueueManagerProps extends QueueProps {
    /**
     * OG事件列表
     */
    ogEvents: string[];
}
interface IEvent extends ActionInfo {
    inner_status?: string;
}

/**
 * 队列原子操作
 */

declare class Queue {
    protected maxLength: any;
    protected userActionSetId: any;
    protected lostActionMaps: Record<string, number>;
    private stack;
    private localStorageKey;
    private timeStamp;
    constructor({ userActionSetId, maxLength }: QueueProps);
    getItems(): IEvent[];
    getStorage(): IEvent[];
    /**
     * 上报丢失条数到日志系统
     */
    reportLostNum(): void;
    /**
     * 获取丢失条数map
     * @returns 丢失条数map
     */
    getLostMaps(): Record<string, number>;
    protected init(): void;
    protected addItem(eventData: IEvent): void;
    protected removeItems(actionIds: string[]): void;
    protected updateForReportFail(actionIds: string[]): void;
    protected updateForReporting(actionIds: string[]): void;
    protected updateAllStack(stack: IEvent[]): void;
    protected updateToStorage(): void;
    /**
     * 更新丢失条数
     * @param sessionId 当前被丢弃行为的会话id
     */
    protected updateLostAction(sessionId: string): void;
    /**
     * 设置时间戳
     */
    protected setTimeStamp(): void;
}

declare class QueueManage extends Queue {
    protected ogEvents: string[];
    constructor({ userActionSetId, maxLength, ogEvents }: QueueManagerProps);
    /**
     * 获取可上报的队列
     * @param reportMaxLength 需要上报的队列最大长度
     */
    getReportableActions(reportMaxLength?: number): IEvent[];
    /**
     * 添加行为到队列
     * @param eventData 单条行为数据
     */
    addAction(eventData: IEvent): TrackResponse;
    /**
     * 根据acition_id删除行为
     * @param actionIds 行为id列表
     */
    removeActions(actionIds: string[]): void;
    /**
     * 上报失败时更新队列状态
     * @param actionIds
     */
    updateActionsForReportFail(actionIds: string[]): void;
    /**
     * 上报时更新队列状态为reporting
     * @param actionIds 需要更新的行为id数组
     */
    updateActionsForReporting(actionIds: string[]): void;
    /**
     * 获取可上报队列长度
     * @returns 可上报队列长度
     */
    getReportableActionsLength(): number;
    /**
     * 按照优先级从大到小给队列排序
     * @param currEvent 当前事件
     * @param stack 当前队列
     */
    sortQueue(currEvent: IEvent, stack: IEvent[]): IEvent[];
    /**
     * 计算优先级
     * @param currTime 当前时间
     * @param event 行为事件
     * @returns 权重值
     */
    private caculateWeight;
    /**
     *
     * @param currTime 当前时间
     * @param actionTime 队列行为时间
     * @returns 权重值
     */
    private formatWeight;
}

interface ConfigManager {
    getRealTimeActionList: () => Array<string>;
    getRequestTimeout: () => number;
}
interface ActionReporterProps {
    getBaseInfo: () => TrackBaseInfo;
    reportLog: (param: LogParams) => void;
    queueManager: QueueManage;
    configManager?: ConfigManager;
    onReportComplete?: Function | undefined;
    onReportFail?: Function | undefined;
}
/**
 * 行为上报类
 */
declare class ActionReporter {
    private static currentRequestCount;
    private static requestConcurrency;
    static setRequestConcurrency(n: number): void;
    private getBaseInfo;
    private reportLog;
    private queueManager;
    private configManager;
    private onReportComplete;
    private onReportFail;
    private cgiBatchSize;
    private reportThreshold;
    private reportDelay;
    private triggerExecuteSend;
    private inspectDelay;
    private inspectTimer;
    private isNeedContinueSend;
    constructor(props: ActionReporterProps);
    /**
     * 批量上报行为，track的时候触发
     */
    batchSend(): void;
    /**
     * 全量上报， flush的时候触发
     */
    flushSend(): void;
    /**
     * 执行指定数量的行为上报逻辑，超过接口batch上限则并行发送请求
     */
    private executeSend;
    /**
     * 组装上报用的行为数据
     */
    private generateActionReportParams;
    /**
     * 上报接口请求成功逻辑处理
     * @param resData
     * @param actionIdList
     */
    private dealSuccessData;
    /**
     * 上报接口请求失败逻辑处理
     * @param errData
     * @param actionIdList
     */
    private dealFailData;
    /**
     * wx.request文档：https://developers.weixin.qq.com/minigame/dev/api/network/request/wx.request.html
     * @param param.data:接口请求参数、param.actionIdList:行为id列表
     * @returns
     */
    private report;
    /**
     * 启动上报定时器
     */
    private startInspectTimer;
    /**
     * 上报回调数据
     */
    private doTrackCallbackFn;
}

declare abstract class SdkExtCommon {
    protected special_method_symbol: symbol;
    /**
     * 付费
     */
    onPurchase(value: number): TrackResponse;
    /**
     * 进入前台
     */
    onEnterForeground(): TrackResponse;
    /**
     * 进入后台
     */
    onEnterBackground(): TrackResponse;
    /**
     * 启动
     */
    onAppStart(): TrackResponse;
    /**
     * 退出
     */
    onAppQuit(): TrackResponse;
    /**
     * 收藏
     */
    onAddToWishlist(): TrackResponse;
    protected wrapTrack(action_type: ActionInfo["action_type"], action_param?: ActionInfo["action_param"]): TrackResponse;
    abstract track(action_type: ActionInfo["action_type"], action_param?: ActionInfo["action_param"]): TrackResponse;
}

declare abstract class sdkExt extends SdkExtCommon {
    /**
     * 注册
     */
    onRegister(): TrackResponse;
    /**
     * 创角
     */
    onCreateRole(name?: string): TrackResponse;
    /**
     * 完成新手教程
     */
    onTutorialFinish(): TrackResponse;
}

declare const initializedInstanceSymbol: unique symbol;
declare class SDK extends sdkExt {
    private static [initializedInstanceSymbol];
    static setRequestConcurrency(n: number): void;
    static setDebug(debug: boolean): void;
    protected env: string;
    protected sdk_version: string;
    protected sdk_name: string;
    protected config?: SdkConfig;
    protected deviceInfo: Partial<DeviceInfo>;
    protected gameInfo: Partial<GameInfo>;
    protected session_id: string;
    protected log_id: number;
    protected queueManage?: QueueManage;
    protected actionReporter?: ActionReporter;
    protected openid?: string;
    protected unionid?: string;
    protected user_unique_id?: string;
    protected inited: boolean;
    protected initErrMsg: string;
    protected onReportComplete?: Function | undefined;
    protected onReportFail?: Function | undefined;
    constructor(sdkConfig: SdkConfig);
    /**
     * 返回初始化状态
     * @returns
     */
    getInitResult(): {
        inited: boolean;
        initErrMsg: string;
    };
    /**
     * 数据上报方法
     * 0: 成功
     * 100: 未完成初始化或重复初始化导致初始化失败，请先完成初始化
     * 101: action_type 或 action_param 参数错误
     * 102: action_param 参数过大，不能超过 xxxx 字符
     * 103: 缓存队列已满，主动丢弃行为
     *
     */
    track(action_type: ActionInfo['action_type'], action_param?: ActionInfo['action_param']): TrackResponse;
    /**
     * 立即上报数据
     */
    flush(): void;
    /**
     * 设置 openId
     * 0: 成功
     * 101: openid 格式错误
     */
    setOpenId(openid: string): {
        code: number;
        message: string;
    };
    /**
       * 设置 unionid
       * 0: 成功
       * 101: unionid 格式错误
       */
    setUnionId(unionid: string): {
        code: number;
        message: string;
    };
    /**
       * 设置 user_unique_id
       * @param user_unique_id
       */
    setUserUniqueId(user_unique_id: string): {
        code: number;
        message: string;
    };
    protected doReportOnEnterBackground(): void;
    /**
       * 获取数据上报时所需要的 info 参数
       */
    private getTrackBaseInfo;
    /**
     * 生成上报队列需要的行为数据
     * @param actionParams 行为参数
     * @param is_sdk_auto_track 是否自动采集
     * @returns
     */
    private createAction;
    /**
     * 给指定行为的inner_param字段追加字段
     * @param action 待追加字段的行为
     * @param key 待添加到inner_param字段的key
     * @param value key对应的value
     */
    private addActionInnerParam;
    /**
     * 不同行为携带的渠道来源判断逻辑
     * @param actionType
     * @returns
     */
    private getChannelByActionType;
    /**
     * 上报日志
     * @param params 日志参数
     */
    private reportLog;
    private useAutoTrack;
    private getAutoProxyRemoteConfig;
    /**
     * 设置openid时，如果值合法缓存到Storage
     */
    private saveValidOpenidToStorage;
}

export { SDK };

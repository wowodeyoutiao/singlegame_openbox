using DG.Tweening;
using UnityEngine;

namespace JEngine.Core
{
    public partial class Updater : MonoBehaviour
    {
        [SerializeField] [Tooltip("热更资源下载地址")] private string resourceUrl = "http://192.168.13.75/yongsitest";

        [SerializeField] [Tooltip("热更资源下载备用地址")]
        private string fallbackUrl = "http://192.168.13.75/yongsitest";

        [SerializeField] private string gameScene = "Assets/HotUpdateResources/Main/Scene/Login.unity";
        [SerializeField] private string mainPackageName = "Main";

        [Tooltip("Simulate是开发模式，Standalone是离线模式，Remote是真机模式")] [SerializeField]
        private UpdateMode mode = UpdateMode.Simulate;

        /// <summary>
        /// WebGL版本号
        /// </summary>
        [SerializeField] private string webVersion = "1.0";


        /// <summary>
        /// 单例
        /// </summary>
        public static Updater Instance => _instance;

        /// <summary>
        /// 主包
        /// </summary>
        public static string MainPackageName => _instance.mainPackageName;

        /// <summary>
        /// 模式
        /// </summary>
        public static UpdateMode Mode => _instance.mode;

        /// <summary>
        /// 热更资源下载地址
        /// </summary>
        public static string ResourceUrl => _instance.resourceUrl;

        /// <summary>
        /// 热更资源下载备用地址
        /// </summary>
        public static string FallbackUrl => _instance.fallbackUrl;

        /// <summary>
        /// WebGL版本号
        /// </summary>
        public static string WebVersion => _instance.webVersion;

        /// <summary>
        /// 更新模式
        /// </summary>
        public enum UpdateMode : byte
        {
            Simulate = 0,
            Standalone = 1,
            Remote = 2,
            WebGL = 3,
            WebGLRemote = 4,
        }

        /// <summary>
        /// 单例
        /// </summary>
        private static Updater _instance;

        /// <summary>
        /// 更新配置
        /// </summary>
        private void Awake()
        {
            if (_instance != null)
            {
                Destroy(_instance.gameObject);
            }
            _instance = this;
            DontDestroyOnLoad(gameObject);
        }
        private void OnEnable()
        {
            AssetMgr.Init();

            StartUpdate();

        }
        /// <summary>
        /// 给按钮拖拽赋值的点击事件，下载更新，用于初始化主包
        /// </summary>
        public void StartUpdate()
        {
            var updater = FindObjectOfType<UpdateScreen>();
            updater.sceneName = gameScene;
            _ = AssetMgr.UpdatePackage(mainPackageName, updater);
        }

        private void OnDestroy()
        {
            MessageBox.Dispose();
        }

        public static void Quit()
        {
#if UNITY_EDITOR
            UnityEditor.EditorApplication.isPlaying = false;
#elif WEIXINMINIGAME
            SDKCtr.Instance.QuitGame();
#else
            Application.Quit();
#endif
        }
    }
}
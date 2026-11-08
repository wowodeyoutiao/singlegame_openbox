using Cysharp.Threading.Tasks;
using JEngine.Core;
using JEngine.Singleton;
using System;
using System.Collections.Generic;
using System.IO;
using UnityEngine;
using UnityEngine.U2D;
using UnityEngine.UI;
//using System.Threading.Tasks;
using YooAsset;

namespace JEngine.Helper
{
    public class ResouseHelper : SingletonTemplate<ResouseHelper>
    {
        //相对路径
        static protected string BaseResousePath = "Assets/HotUpdateResources/";

        private Dictionary<string, SpriteAtlas> spriteAtlasDict = null;
        private Dictionary<string, Sprite> spriteDict = null;
        public ResouseHelper()
        {
            spriteAtlasDict = new Dictionary<string, SpriteAtlas>();
            spriteDict = new Dictionary<string, Sprite>();
        }

        public void Clear()
        {
            spriteAtlasDict.Clear();
        }

        public Material LoadMaterial(string abname, string packge = "Main")
        {
            abname = BaseResousePath + abname;
            if (!abname.Contains(".mat"))
            {
                abname = abname + ".mat";
            }
            return AssetMgr.Load<Material>(abname, packge);
        }

        public async UniTask LoadMaterialAsync(string path, Action<Material> action, string packge = "Main")
        {
            path = BaseResousePath + packge + "/" + path;
            (var obj, var handle) = await AssetMgr.LoadAsyncWithHandle<Material>(path, packge);
            await UniTask.SwitchToMainThread();
            action?.Invoke(obj);
            handle.Release();
        }

        public async UniTask LoadGameObj(string path, Action<GameObject> ObjAction, string packge = "Main")
        {
            if (!path.Contains(".prefab"))
                path = BaseResousePath + packge + "/" + path + ".prefab";

            (var obj, var handle) = await AssetMgr.LoadAsyncWithHandle<GameObject>(path, packge);
            ObjAction?.Invoke(obj);
            //handle.Release();
        }


        public Texture2D LoadTexture(string abname, string packge = "Main")
        {
            return AssetMgr.Load<Texture2D>(abname, packge);
        }

        public async UniTask LoadTextAssetAsync(string path, Action<string> action, string packge = "Main")
        {
            path = BaseResousePath + packge + "/" + path;
            try
            {
                (var obj, var handle) = await AssetMgr.LoadAsyncWithHandle<TextAsset>(path, packge);
                //var obj = await AssetMgr.LoadAsync<TextAsset>(path, packge);
                await UniTask.SwitchToMainThread(); // 确保在主线程执行回调
                if (obj == null)
                {
                    UnityEngine.Debug.LogError("[CFG] TextAsset missing: " + path);
                    action?.Invoke("[]"); // keep caller's progress moving instead of dead-locking
                }
                else
                {
                    action?.Invoke(obj.text);
                }
                handle.Release();
            }
            catch (System.Exception e)
            {
                UnityEngine.Debug.LogError("[CFG] LoadTextAssetAsync error: " + path + " | " + e);
            }
        }


        public T LoadAsset<T>(string path, string packge = "Main") where T : UnityEngine.Object
        {
            path = BaseResousePath + packge + "/" + path;
            return AssetMgr.Load<T>(path, packge);
        }

        public async UniTask LoadAudioClipAssetAsync(string path, Action<AudioClip> action, string packge = "Main")
        {
            path = BaseResousePath + packge + "/" + path;
            //var obj = await AssetMgr.LoadAsync<AudioClip>(path, packge);
            (var obj, var handle) = await AssetMgr.LoadAsyncWithHandle<AudioClip>(path, packge);
            action?.Invoke(obj);
            //handle.Release();
        }


        //public async Task<T> LoadAssetAsync<T>(string path, string packge = "Main") where T : UnityEngine.Object
        //{
        //    path = BaseResousePath + packge + "/" + path;
        //    var textAsset = await AssetMgr.LoadAsync<T>(path, packge);
        //    return textAsset;
        //}


        public async UniTask LoadSprite(string path, Image img, string packge = "Main")
        {
            path = BaseResousePath + packge + "/" + path;
            if (spriteDict.ContainsKey(path))
            {
                img.sprite = spriteDict[path];
                return;
            }
            (var sprite, var handle) = await AssetMgr.LoadAsyncWithHandle<Sprite>(path, packge);
            //img.sprite = await AssetMgr.LoadAsync<Sprite>(path, packge);
            img.sprite = sprite;
            spriteDict[path] = sprite;
            //handle.Release();
        }

        public async UniTask LoadSprite(string path, SpriteRenderer img, string packge = "Main")
        {
            path = BaseResousePath + packge + "/" + path;
            if (spriteDict.ContainsKey(path))
            {
                img.sprite = spriteDict[path];
                return;
            }
            (var sprite, var handle) = await AssetMgr.LoadAsyncWithHandle<Sprite>(path, packge);
            //img.sprite = await AssetMgr.LoadAsync<Sprite>(path, packge);
            img.sprite = sprite;
            spriteDict[path] = sprite;
            //handle.Release();
        }

        public async void LoadSpriteByAtlas(string atlas, string iconName, Image img, string packge = "Main", Action callBack = null)
        {
            if (atlas == "" || atlas == null || iconName == "")
            {
                img.sprite = null;
                return; //原实现缺 return，会继续去加载空路径导致异常/白图
            }
            var _atls = BaseResousePath + packge + "/" + atlas + "/_Atlas";
            if (!_atls.Contains(".spriteatlas"))
            {
                _atls = _atls + ".spriteatlas";
            }
            AssetHandle handle = null;
            if (!spriteAtlasDict.ContainsKey(atlas))
            {
                (var sprite, handle) = await AssetMgr.LoadAsyncWithHandle<SpriteAtlas>(_atls, packge);
                if (!spriteAtlasDict.ContainsKey(atlas))
                {
                    spriteAtlasDict.Add(atlas, sprite);
                }
            }
            img.sprite = spriteAtlasDict[atlas].GetSprite(iconName);
            callBack?.Invoke();
            //if (handle != null)
            //    handle.Release();
        }

        public async UniTaskVoid LoadSpriteByAtlas(string atlas, string iconName, SpriteRenderer img, string packge = "Main")
        {
            if (atlas == "" || atlas == null || iconName == "")
            {
                img.sprite = null;
            }
            var _atls = BaseResousePath + packge + "/" + atlas + "/_Atlas";
            if (!_atls.Contains(".spriteatlas"))
            {
                _atls = _atls + ".spriteatlas";
            }
            AssetHandle handle = null;
            if (!spriteAtlasDict.ContainsKey(atlas))
            {
                (var sprite, handle) = await AssetMgr.LoadAsyncWithHandle<SpriteAtlas>(_atls, packge);
                //var sprite = await AssetMgr.LoadAsync<SpriteAtlas>(_atls, packge);
                if (!spriteAtlasDict.ContainsKey(atlas))
                {
                    spriteAtlasDict.Add(atlas, sprite);
                }
            }
            img.sprite = spriteAtlasDict[atlas].GetSprite(iconName);
            if (handle != null)
                handle.Release();
        }


        public FileStream loadFileStream(string path, string packge = "Main")
        {
            path = BaseResousePath + packge + "/" + path;

            if (!File.Exists(path))
                return null;

            FileStream filestream = new FileStream(path, FileMode.Open);

            return filestream;
        }
    }
}

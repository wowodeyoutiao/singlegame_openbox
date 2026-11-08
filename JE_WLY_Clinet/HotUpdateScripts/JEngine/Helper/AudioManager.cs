using Cysharp.Threading.Tasks;
using HotUpdateScripts;
using JEngine.Helper;
using System;
using System.Collections.Generic;
using UnityEngine;

public class AudioManager : MonoBehaviour
{

    private bool _isOpenSoundEffect = true;
    public bool isOpenSoundEffect
    {
        get
        {
            return _isOpenSoundEffect;
        }
        set
        {
            if (_isOpenSoundEffect != value)
            {
                GameUtil.Instance.SetPlayerPrefs(GameConst.SystemSettingOption.SoundEffect.ToString(), value ? 1 : 0);
                _isOpenSoundEffect = value;
            }
        }
    }

    private bool _isOpenBgMusic = true;
    public bool isOpenBgMusic
    {
        get
        {
            return _isOpenBgMusic;
        }
        set
        {
            if (value != _isOpenBgMusic)
            {
                _isOpenBgMusic = value;
                GameUtil.Instance.SetPlayerPrefs(GameConst.SystemSettingOption.BgMusic.ToString(), value ? 1 : 0);
            }
        }
    }

    private string BasePath = "Other/Audios/";


    string prevBgMusic;
    string prevLoopSound;

    private AudioSource m_bgMusicSource;
    private AudioSource m_shotMusicSource;
    private AudioSource m_LoopSoundSource;

    private Dictionary<string, AudioClip> m_cacheAudio = new Dictionary<string, AudioClip>();

    [Header("同时播放的音效上限")]
    public int maxConcurrentSE = 10;

    private static AudioManager instance;
    public static AudioManager Instance
    {
        get
        {
            if (instance == null)
            {
                GameObject audioGo = new GameObject("AudioManager");
                GameObject shotSource = new GameObject("OneShot");
                GameObject loopSource = new GameObject("LoopAudio");
                shotSource.transform.SetParent(audioGo.transform);

                instance = audioGo.AddComponent<AudioManager>();
                instance.m_bgMusicSource = audioGo.AddComponent<AudioSource>();
                instance.m_shotMusicSource = shotSource.AddComponent<AudioSource>();
                instance.m_LoopSoundSource = loopSource.AddComponent<AudioSource>();

                instance._isOpenBgMusic = Convert.ToBoolean(GameUtil.Instance.GetPlayerPrefsInt(GameConst.SystemSettingOption.BgMusic.ToString(), 1));
                instance._isOpenSoundEffect = Convert.ToBoolean(GameUtil.Instance.GetPlayerPrefsInt(GameConst.SystemSettingOption.SoundEffect.ToString(), 1));

                DontDestroyOnLoad(audioGo);
            }
            return instance;
        }
    }

    void Awake()
    {

    }


    public void SetBgAudioVolume(float value)
    {
        m_bgMusicSource.volume = value;
    }

    public float GetBgAudioVolume()
    {
        return m_bgMusicSource.volume;
    }

    public void SetShotAudioVolume(float value)
    {
        m_shotMusicSource.volume = value;
        m_LoopSoundSource.volume = value;
    }

    public float GetShotAudioVolume()
    {
        return m_shotMusicSource.volume;
    }

    /// <summary>
    /// 通过配置播放音效
    /// </summary>
    /// <param name="_sound"></param>
    public void PlaySound(AudioKey.Sound _sound)
    {
        //Debug.LogError(_sound.ToString());
        PlaySound(AudioKey.SoundDic[_sound]);
    }

    /// <summary>
    /// 通过相对路径播放音效
    /// </summary>
    /// <param name="clipPath"></param>
    /// <param name="volumeScale"></param>
    public async void PlaySound(string clipPath, float volumeScale = 1f)
    {
        if (isOpenSoundEffect)
        {
            if (m_cacheAudio.ContainsKey(clipPath))
            {
                PlaySound(m_cacheAudio[clipPath], volumeScale);
            }
            else
            {
                await ResouseHelper.Instance.LoadAudioClipAssetAsync(BasePath + clipPath, (clip) =>
                {
                    m_cacheAudio[clipPath] = clip;
                    PlaySound(clip, volumeScale);
                });
            }
        }
    }

    private int NowPlayCount = 0;
    private void PlaySound(AudioClip clip, float volumeScale = 1f)
    {

        if (isOpenSoundEffect && NowPlayCount <= maxConcurrentSE)
        {
            if (clip == null) return;
            m_shotMusicSource.PlayOneShot(clip, volumeScale);
            NowPlayCount++;
            UniCoroutinePlaySoundCount(clip.length).Forget();
        }
    }

    private async UniTask UniCoroutinePlaySoundCount(float time)
    {
        await UniTask.Delay((int)(time * 1000));
        NowPlayCount--;
    }

    /// <summary>
    /// 通过配置播放背景音乐
    /// </summary>
    /// <param name="_sound"></param>
    public void PlayBgMusic(AudioKey.Sound _sound)
    {
        PlayBgMusic(AudioKey.SoundDic[_sound]);
    }

    /// <summary>
    /// 通过路径播放背景音乐
    /// </summary>
    /// <param name="clipPath"></param>
    public async void PlayBgMusic(string clipPath)
    {
        if (prevBgMusic == clipPath)
        {
            return;
        }
        if (m_cacheAudio.ContainsKey(clipPath))
        {
            PlayBgMusic(m_cacheAudio[clipPath]);
        }
        else
        {
            await ResouseHelper.Instance.LoadAudioClipAssetAsync(BasePath + clipPath, (clip) =>
            {
                PlayBgMusic(clip);
            });
        }

    }

    private void PlayBgMusic(AudioClip clip)
    {
        if (clip == null)
        {
            return;
        }
        string clipName = clip.name;

        if (prevBgMusic == clipName) { return; } else { StopBgMusic(); }


        m_bgMusicSource.loop = true;
        m_bgMusicSource.clip = clip;

        if (m_bgMusicSource.clip != null && m_bgMusicSource.isPlaying == false && isOpenBgMusic)
        {
            m_bgMusicSource.Play();
        }

        prevBgMusic = clipName;
    }

    public void StopBgMusic()
    {
        if (m_bgMusicSource != null && m_bgMusicSource.isPlaying)
        {
            m_bgMusicSource.Stop();
            m_bgMusicSource.clip = null;

        }
        prevBgMusic = null;
    }

    public void ReSumeBgMusic()
    {
        if (isOpenBgMusic && m_bgMusicSource != null && m_bgMusicSource.clip != null)
        {
            m_bgMusicSource.Play();
        }
    }


    public void PauseBgMusic()
    {
        if (m_bgMusicSource != null && m_bgMusicSource.clip != null)
        {
            m_bgMusicSource.Pause();
        }
    }

    /// <summary>
    /// 通过配置播放循环音效
    /// </summary>
    /// <param name="_sound"></param>
    public void PlayLoopSound(AudioKey.Sound _sound)
    {
        PlayLoopSound(AudioKey.SoundDic[_sound]);
    }

    /// <summary>
    /// 通过路径播放循环音效
    /// </summary>
    /// <param name="clipPath"></param>
    /// <param name="volumeScale"></param>
    public async void PlayLoopSound(string clipPath, float volumeScale = 1f)
    {
        if (prevLoopSound == clipPath)
            return;
        if (m_cacheAudio.ContainsKey(clipPath))
        {
            PlayLoopSound(m_cacheAudio[clipPath], volumeScale);
        }
        else
        {
            await ResouseHelper.Instance.LoadAudioClipAssetAsync(BasePath + clipPath, (clip) =>
            {
                m_cacheAudio[clipPath] = clip;
                PlayLoopSound(clip, volumeScale);
            });
            // AudioClip clip = await ResouseHelper.Instance.LoadAssetAsync<AudioClip>(BasePath + clipPath+".ogg");
        }

    }

    private void PlayLoopSound(AudioClip clip, float volumeScale = 1f)
    {
        if (clip == null)
        {
            return;
        }
        string clipName = clip.name;
        if (prevLoopSound == clipName) { return; } else { StopLoopSound(); }
        m_LoopSoundSource.loop = true;
        m_LoopSoundSource.clip = clip;
        if (m_LoopSoundSource.clip != null && m_LoopSoundSource.isPlaying == false)
        {
            m_LoopSoundSource.Play();
        }
        prevBgMusic = clipName;
    }

    public void StopLoopSound()
    {
        if (m_LoopSoundSource != null && m_LoopSoundSource.isPlaying)
        {
            m_LoopSoundSource.Stop();
            m_LoopSoundSource.clip = null;
        }
        prevLoopSound = null;
    }


    public void OnMusic(bool isOn)
    {
        isOpenBgMusic = isOn;
        if (isOpenBgMusic == false)
        {
            PauseBgMusic();
        }
        else
        {
            if (m_bgMusicSource != null && m_bgMusicSource.isPlaying == false)
            {
                m_bgMusicSource.loop = true;
                if (m_bgMusicSource.clip != null)
                {
                    m_bgMusicSource.Play();
                }
            }
        }

    }
    public void OnSound(bool isOn)
    {
        isOpenSoundEffect = isOn;
    }



}

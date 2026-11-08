using JEngine.Singleton;
using Spine.Unity;
using System;
using System.Collections.Generic;
using UnityEngine;
using static Spine.AnimationState;

namespace JEngine.Helper
{
    public class SpineUtil : SingletonTemplate<SpineUtil>
    {
        private readonly Dictionary<SkeletonAnimation, TrackEntryDelegate> CompleteSkeAnimDic = new Dictionary<SkeletonAnimation, TrackEntryDelegate>();
        private readonly Dictionary<SkeletonGraphic, TrackEntryDelegate> CompleteSkeletonGraphicDic = new Dictionary<SkeletonGraphic, TrackEntryDelegate>();

        private static readonly Dictionary<int, string> KuangSpineName = new()
        {
            [2] = "lv",
            [3] = "lan",
            [4] = "zi",
            [6] = "cheng",
            [7] = "hong"

        };
        public void PlayOnce(SkeletonGraphic skeletonGraphic, string animationName, Action callBack, float timeScale = 1)
        {
            CheakAndRemove(skeletonGraphic);
            skeletonGraphic.AnimationState.ClearTracks();
            if (!string.IsNullOrEmpty(animationName))
                skeletonGraphic.AnimationState.SetAnimation(0, animationName, false);
            skeletonGraphic.AnimationState.TimeScale = timeScale;

            TrackEntryDelegate compCallback = _ =>
            {
                CheakAndRemove(skeletonGraphic);
                callBack?.Invoke();
            };

            skeletonGraphic.AnimationState.Complete += compCallback;
            CompleteSkeletonGraphicDic[skeletonGraphic] = compCallback;
        }

        public void PlayOnce(SkeletonAnimation skeletonAnimation, string animationName, Action callBck, float timeScale = 1)
        {
            CheakAndRemove(skeletonAnimation);
            skeletonAnimation.state.ClearTracks();
            if (!string.IsNullOrEmpty(animationName))
                skeletonAnimation.state.SetAnimation(0, animationName, false);
            skeletonAnimation.state.TimeScale = timeScale;

            TrackEntryDelegate compCallback = _ =>
            {
                CheakAndRemove(skeletonAnimation);
                callBck?.Invoke();
            };
            skeletonAnimation.state.Complete += compCallback;
            CompleteSkeAnimDic[skeletonAnimation] = compCallback;
        }

        public void PlayOnce(Transform root, string animationName, Action callBck, float timeScale = 1)
        {
            if (root.GetComponent<SkeletonGraphic>() != null)
            {
                PlayOnce(root.GetComponent<SkeletonGraphic>(), animationName, callBck, timeScale);
            }
            else if (root.GetComponent<SkeletonAnimation>() != null)
            {
                PlayOnce(root.GetComponent<SkeletonAnimation>(), animationName, callBck, timeScale);
            }
            else
            {
                UnityEngine.Debug.LogError("PlayOnce 没有找到SkeletonGraphic或SkeletonAnimation组件");
            }

        }

        public void PlayLoop(Transform root, string animationName, float timeScale = 1)
        {
            if (root.GetComponent<SkeletonGraphic>() != null)
            {
                PlayLoop(root.GetComponent<SkeletonGraphic>(), animationName, timeScale);
            }
            else if (root.GetComponent<SkeletonAnimation>() != null)
            {
                PlayLoop(root.GetComponent<SkeletonAnimation>(), animationName, timeScale);
            }
            else
            {
                UnityEngine.Debug.LogError("PlayLoop 没有找到SkeletonGraphic或SkeletonAnimation组件");
            }
        }

        public void PlayItemSpineLoop(Transform root, int quality, float timeScale = 1)
        {
            var animationName = KuangSpineName[4];
            if (KuangSpineName.TryGetValue(quality, out var name))
                animationName = name;

            PlayLoop(root, animationName, timeScale);
        }

        public void PlayLoop(SkeletonGraphic skeletonGraphic, string animationName, float timeScale = 1)
        {
            CheakAndRemove(skeletonGraphic);
            skeletonGraphic.Skeleton.SetToSetupPose();
            skeletonGraphic.AnimationState.ClearTracks();
            if (!string.IsNullOrEmpty(animationName))
                skeletonGraphic.AnimationState.SetAnimation(0, animationName, true);

            skeletonGraphic.AnimationState.TimeScale = timeScale;
        }

        public void PlayLoop(SkeletonAnimation skeletonAnimation, string animationName, float timeScale = 1)
        {
            CheakAndRemove(skeletonAnimation);
            skeletonAnimation.skeleton.SetToSetupPose();
            skeletonAnimation.state.ClearTracks();
            if (!string.IsNullOrEmpty(animationName))
                skeletonAnimation.state.SetAnimation(0, animationName, true);

            skeletonAnimation.state.TimeScale = timeScale;
        }

        public void PlayQueue(Transform root, List<string> AnimNameQueue, Action callBck, bool isLoop = false)
        {
            if (root.GetComponent<SkeletonGraphic>() != null)
            {
                PlayQueue(root.GetComponent<SkeletonGraphic>(), AnimNameQueue, callBck);
            }
            else if (root.GetComponent<SkeletonAnimation>() != null)
            {
                PlayQueue(root.GetComponent<SkeletonAnimation>(), AnimNameQueue, callBck, isLoop);
            }
            else
            {
                UnityEngine.Debug.LogError("PlayLoop 没有找到SkeletonGraphic或SkeletonAnimation组件");
            }
        }

        public void PlayQueue(SkeletonGraphic skeletonGraphic, List<string> AnimNameQueue, Action callBck, bool isLoop = false)
        {
            CheakAndRemove(skeletonGraphic);
            skeletonGraphic.AnimationState.ClearTracks();

            skeletonGraphic.AnimationState.SetAnimation(0, AnimNameQueue[0], AnimNameQueue.Count == 1 ? isLoop : false);
            AnimNameQueue.RemoveAt(0);

            TrackEntryDelegate compCallback = _ =>
            {
                CheakAndRemove(skeletonGraphic);
                int QueueCount = AnimNameQueue.Count;
                if (QueueCount > 0)
                {
                    PlayQueue(skeletonGraphic, AnimNameQueue, callBck, isLoop);
                }
                if (QueueCount == 1)
                    callBck?.Invoke();

            };
            skeletonGraphic.AnimationState.Complete += compCallback;
            CompleteSkeletonGraphicDic[skeletonGraphic] = compCallback;
        }

        public void PlayQueue(SkeletonAnimation skeletonAnimation, List<string> AnimNameQueue, Action callBck, bool isLoop = false)
        {
            CheakAndRemove(skeletonAnimation);
            skeletonAnimation.skeleton.SetToSetupPose();
            skeletonAnimation.state.ClearTracks();
            skeletonAnimation.state.SetAnimation(0, AnimNameQueue[0], AnimNameQueue.Count == 1 ? isLoop : false);
            AnimNameQueue.RemoveAt(0);

            TrackEntryDelegate compCallback = _ =>
            {
                CheakAndRemove(skeletonAnimation);
                int QueueCount = AnimNameQueue.Count;
                if (QueueCount > 0)
                {
                    PlayQueue(skeletonAnimation, AnimNameQueue, callBck, isLoop);
                }
                if (QueueCount == 1)
                    callBck?.Invoke();
            };
            skeletonAnimation.state.Complete += compCallback;
            CompleteSkeAnimDic[skeletonAnimation] = compCallback;
        }

        private void CheakAndRemove(SkeletonAnimation skeletonAnimation)
        {
            if (CompleteSkeAnimDic.ContainsKey(skeletonAnimation))
            {
                skeletonAnimation.state.Complete -= CompleteSkeAnimDic[skeletonAnimation];
                CompleteSkeAnimDic.Remove(skeletonAnimation);
            }
        }

        private void CheakAndRemove(SkeletonGraphic skeletonAnimation)
        {
            if (CompleteSkeletonGraphicDic.ContainsKey(skeletonAnimation))
            {
                skeletonAnimation.AnimationState.Complete -= CompleteSkeletonGraphicDic[skeletonAnimation];
                CompleteSkeletonGraphicDic.Remove(skeletonAnimation);
            }
        }
    }
}

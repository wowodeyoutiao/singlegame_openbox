using Spine.Unity;
using UnityEngine;

public class SpineCtr : MonoBehaviour
{
    public float Time;
    public string AnimationName;

    private SkeletonGraphic skeletonGraphic;
    private SkeletonAnimation skeletonAnimation;
    private void Awake()
    {
        if (GetComponent<SkeletonGraphic>() != null)
        {
            skeletonGraphic = GetComponent<SkeletonGraphic>();
        }

        if (GetComponent<SkeletonAnimation>() != null)
        {
            skeletonAnimation = GetComponent<SkeletonAnimation>();
        }
    }

    void OnEnable()
    {
        Invoke("PlaySpine", Time);
    }

    private void PlaySpine()
    {
        if (skeletonGraphic != null)
        {
            skeletonGraphic.AnimationState.ClearTracks();
            skeletonGraphic.AnimationState.SetAnimation(0, AnimationName, false);
        }

        if (skeletonAnimation != null)
        {
            skeletonAnimation.AnimationState.ClearTracks();
            skeletonAnimation.AnimationState.SetAnimation(0, AnimationName, false);
        }
    }
}

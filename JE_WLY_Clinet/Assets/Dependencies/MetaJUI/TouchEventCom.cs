using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.Events;
using UnityEngine.EventSystems;
using UnityEngine.UI;


[System.Serializable]
public class TouchEvent : UnityEvent<PointerEventData>  
{
}

public class TouchEventCom : EventTrigger  
{
    private TouchEvent clickBeginEvent;
    private TouchEvent clickEndEvent;
    private TouchEvent clickEvent;
    private TouchEvent dragBeginEvent;
    private TouchEvent dragEndEvent;
    private TouchEvent dragEvent;
    private TouchEvent pointerUpEvent;
    private TouchEvent pointerEnterEvent;
    private TouchEvent pointerDownEvent;
    private TouchEvent pointerExitEvent;

    public ScrollRect scrollRect;

    private void Awake()
    {
        if (!scrollRect)
        {
            scrollRect = GetComponentInParent<ScrollRect>();
        }
    }

    void Start()
    {
    }

    public static void PassEventPointerDown(PointerEventData data , object function)
    {   
        PassEvent(data, (ExecuteEvents.EventFunction<IPointerDownHandler>)function);
    }

    public static void PassEventPointerUp(PointerEventData data , object function)
    {
        PassEvent(data, (ExecuteEvents.EventFunction<IPointerUpHandler>)function);
    }

    public static void PassEventSubmit(PointerEventData data , object function)
    {
        PassEvent(data, (ExecuteEvents.EventFunction<ISubmitHandler>)function);
    }

    public static void PassEventClick(PointerEventData data , object function)
    {
        PassEvent(data, (ExecuteEvents.EventFunction<IPointerClickHandler>)function);
    }

    public static void PassEvent<T>(PointerEventData data, ExecuteEvents.EventFunction<T> function) where T : IEventSystemHandler
    {
        List<RaycastResult> results = new List<RaycastResult>();
        EventSystem.current.RaycastAll(data, results);
        GameObject current = data.pointerCurrentRaycast.gameObject;
        for (int i = 0; i < results.Count; i++)
        {
            if (current != results[i].gameObject)
            {
                ExecuteEvents.Execute(results[i].gameObject, data, function);
                break;
            }
        }
    }

    public override void OnBeginDrag(PointerEventData data)
    {
        if (scrollRect)
        {
            scrollRect.OnBeginDrag(data);
        }
        if (dragBeginEvent == null)
        {
            return;
        }
        dragBeginEvent.Invoke(data);
    }

    public override void OnDrag(PointerEventData data)
    {
        if (scrollRect)
        {
            scrollRect.OnDrag(data);
        }
        if (dragEvent == null)
        {
            return;
        }
        dragEvent.Invoke(data);
    }

    public override void OnEndDrag(PointerEventData data)
    {
        if (scrollRect)
        {
            scrollRect.OnEndDrag(data);
        }
        if (dragEndEvent == null)
        {
            return;
        }
        dragEndEvent.Invoke(data);
    }

    public override void OnPointerClick(PointerEventData data)
    {
        if (clickEvent == null)
        {
            return;
        }
        clickEvent.Invoke(data);
    }

    public override void OnPointerDown(PointerEventData data)
    {
        if (pointerDownEvent == null)
        {
            return;
        }
        pointerDownEvent.Invoke(data);
    }

    public override void OnPointerEnter(PointerEventData data)
    {
        if (pointerEnterEvent == null)
        {
            return;
        }
        pointerEnterEvent.Invoke(data);
    }

    public override void OnPointerUp(PointerEventData data)
    {
        if (pointerUpEvent == null)
        {
            return;
        }
        pointerUpEvent.Invoke(data);
    }

    public override void OnPointerExit(PointerEventData data)
    {
        if (pointerExitEvent == null)
        {
            return;
        }
        pointerExitEvent.Invoke(data);
    }

    public void AddDragBeginListener(UnityAction<PointerEventData> action)
    {
        if (dragBeginEvent == null)
        {
            dragBeginEvent = new TouchEvent();  
        }
        dragBeginEvent.RemoveAllListeners();
        dragBeginEvent.AddListener(action);
    }

    public void AddDragEndListener(UnityAction<PointerEventData> action)
    {
        if (dragEndEvent == null)
        {
            dragEndEvent = new TouchEvent();
        }
        dragEndEvent.RemoveAllListeners();
        dragEndEvent.AddListener(action);
    }

    public void AddDragListener(UnityAction<PointerEventData> action)
    {
        if (dragEvent == null)
        {
            dragEvent = new TouchEvent();
        }
        dragEvent.RemoveAllListeners();
        dragEvent.AddListener(action);
    }

    public void AddClickBeginListener(UnityAction<PointerEventData> action)
    {
        if (clickBeginEvent == null)
        {
            clickBeginEvent = new TouchEvent();
        }
        clickBeginEvent.RemoveAllListeners();
        clickBeginEvent.AddListener(action);
    }

    public void AddClickEndListener(UnityAction<PointerEventData> action)
    {
        if (clickEndEvent == null)
        {
            clickEndEvent = new TouchEvent();
        }
        clickEndEvent.RemoveAllListeners();
        clickEndEvent.AddListener(action);
    }

    public void AddClickListener(UnityAction<PointerEventData> action)
    {
        if (clickEvent == null)
        {
            clickEvent = new TouchEvent();
        }
        clickEvent.RemoveAllListeners();
        clickEvent.AddListener(action);
    }

    public void AddPointerEnterListener(UnityAction<PointerEventData> action)
    {
        if (pointerEnterEvent == null)
        {
            pointerEnterEvent = new TouchEvent();
        }
        pointerEnterEvent.RemoveAllListeners();
        pointerEnterEvent.AddListener(action);
    }

    public void AddPointerUpListener(UnityAction<PointerEventData> action)
    {
        if (pointerUpEvent == null)
        {
            pointerUpEvent = new TouchEvent();
        }
        pointerUpEvent.RemoveAllListeners();
        pointerUpEvent.AddListener(action);
    }

    public void AddPointerDownListener(UnityAction<PointerEventData> action)
    {
        if (pointerDownEvent == null)
        {
            pointerDownEvent = new TouchEvent();
        }
        pointerDownEvent.RemoveAllListeners();
        pointerDownEvent.AddListener(action);
    }

    public void AddPointerExitListener(UnityAction<PointerEventData> action)
    {
        if (pointerExitEvent == null)
        {
            pointerExitEvent = new TouchEvent();
        }
        pointerExitEvent.RemoveAllListeners();
        pointerExitEvent.AddListener(action);
    }

    public void RemoveAllListeners()
    {
        if(pointerExitEvent != null)
        {
            pointerExitEvent.RemoveAllListeners();
        }
        if(pointerDownEvent != null)
        {
            pointerDownEvent.RemoveAllListeners();
        }
        if(pointerEnterEvent != null)
        {
            pointerEnterEvent.RemoveAllListeners();
        }
        if(pointerUpEvent != null)
        {
            pointerUpEvent.RemoveAllListeners();
        }
        if(dragEvent != null)
        {
            dragEvent.RemoveAllListeners();
        }
        if(dragEndEvent != null)
        {
            dragEndEvent.RemoveAllListeners();
        }
        if(dragBeginEvent != null)
        {
            dragBeginEvent.RemoveAllListeners();
        }
        if(clickEvent != null)
        {
            clickEvent.RemoveAllListeners();
        }
        if(clickEndEvent != null)
        {
            clickEndEvent.RemoveAllListeners();
        }
        if(clickBeginEvent != null)
        {
            clickBeginEvent.RemoveAllListeners();
        }
    }
}

using UnityEngine;
using System.Collections;
using UnityEngine.EventSystems;

public class EventTriggerListener : EventTrigger
{
    public delegate void PointerEventDelegate(PointerEventData eventData);
    public delegate void BaseEventDelegate(BaseEventData eventData);
    public PointerEventDelegate onClick;
    public PointerEventDelegate onDown;
    public PointerEventDelegate onEnter;
    public PointerEventDelegate onExit;
    public PointerEventDelegate onUp;
    public PointerEventDelegate onDrag;
    public PointerEventDelegate onBeginDrag;
    public PointerEventDelegate onEndDrag;
    public BaseEventDelegate onSelect;
    public BaseEventDelegate onUpdateSelect;

    static public EventTriggerListener Get(GameObject obj)
    {
        EventTriggerListener listener = obj.GetComponent<EventTriggerListener>();
        if (listener == null) listener = obj.AddComponent<EventTriggerListener>();
        return listener;
    }
    public override void OnPointerClick(PointerEventData eventData)
    {
        if (onClick != null) onClick(eventData);
    }
    public override void OnDrag(PointerEventData eventData)
    {
        if (onDrag != null) onDrag(eventData);
    }

    public override void OnBeginDrag(PointerEventData eventData)
    {
        if (onBeginDrag != null) onBeginDrag(eventData);
    }

    public override void OnEndDrag(PointerEventData eventData)
    {
        if (onEndDrag != null) onEndDrag(eventData);
    }
    public override void OnPointerDown(PointerEventData eventData)
    {
        if (onDown != null) onDown(eventData);
    }
    public override void OnPointerEnter(PointerEventData eventData)
    {
        if (onEnter != null) onEnter(eventData);
    }
    public override void OnPointerExit(PointerEventData eventData)
    {
        if (onExit != null) onExit(eventData);
    }
    public override void OnPointerUp(PointerEventData eventData)
    {
        if (onUp != null) onUp(eventData);
    }
    public override void OnSelect(BaseEventData eventData)
    {
        if (onSelect != null) onSelect(eventData);
    }
    public override void OnUpdateSelected(BaseEventData eventData)
    {
        if (onUpdateSelect != null) onUpdateSelect(eventData);
    }
}
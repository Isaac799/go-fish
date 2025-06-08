class DragAndDrop {
  constructor(containerId) {
    this.container = document.getElementById(containerId);
    if (!this.container) return;

    this.draggedElement = null;
    this.startPointerX = 0;
    this.startPointerY = 0;
    this.startElemX = 0;
    this.startElemY = 0;

    this.onPointerDown = this.onPointerDown.bind(this);
    this.onPointerMove = this.onPointerMove.bind(this);
    this.onPointerUp = this.onPointerUp.bind(this);

    this.container.addEventListener("pointerdown", this.onPointerDown);
    this.container.addEventListener("pointermove", this.onPointerMove);
    this.container.addEventListener("pointerup", this.onPointerUp);
    this.container.addEventListener("pointercancel", this.onPointerUp);
  }

  getCurrentTranslation(element) {
    const transformStyle = getComputedStyle(element).transform;
    if (!transformStyle || transformStyle === "none") {
      return [0, 0];
    }
    const matrix = new DOMMatrixReadOnly(transformStyle);
    return [matrix.m41, matrix.m42];
  }

  onPointerDown(event) {
    if (!event.target.classList.contains("node")) return;

    this.draggedElement = event.target;
    [this.startElemX, this.startElemY] = this.getCurrentTranslation(
      this.draggedElement
    );
    this.startPointerX = event.clientX;
    this.startPointerY = event.clientY;

    this.draggedElement.style.cursor = "grabbing";
    this.draggedElement.setPointerCapture(event.pointerId);
  }

  onPointerMove(event) {
    if (!this.draggedElement) return;

    const deltaX = event.clientX - this.startPointerX;
    const deltaY = event.clientY - this.startPointerY;

    let newX = this.startElemX + deltaX;
    let newY = this.startElemY + deltaY;

    const containerBounds = this.container.getBoundingClientRect();
    const elementBounds = this.draggedElement.getBoundingClientRect();

    // constrain inside container
    const maxX = containerBounds.width - elementBounds.width;
    const maxY = containerBounds.height - elementBounds.height;
    newX = Math.max(0, Math.min(newX, maxX));
    newY = Math.max(0, Math.min(newY, maxY));

    this.draggedElement.style.transform = `translate(${newX}px, ${newY}px)`;
  }

  onPointerUp(event) {
    if (!this.draggedElement) return;

    const [finalX, finalY] = this.getCurrentTranslation(this.draggedElement);

    // update hidden input values prior to submission
    const elementId = this.draggedElement.id;
    if (elementId) {
      const inputX = this.draggedElement.querySelector(
        `input[name="${elementId}x"]`
      );
      const inputY = this.draggedElement.querySelector(
        `input[name="${elementId}y"]`
      );
      if (inputX) inputX.value = Math.round(finalX);
      if (inputY) inputY.value = Math.round(finalY);
    }

    this.draggedElement.style.cursor = "grab";
    this.draggedElement.releasePointerCapture(event.pointerId);
    this.draggedElement = null;

    // causes hx-trigger to swap the hx-target (container)
    document.body.dispatchEvent(new Event("positionUpdated"));
  }
}

window.addEventListener("DOMContentLoaded", () => {
  new DragAndDrop("container");
  document.body.addEventListener("htmx:afterSwap", () => {
    new DragAndDrop("container");
  });
});

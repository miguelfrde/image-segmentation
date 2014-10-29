function changeImage(activeImage, hiddenImage) {
  $('#show-' + activeImage).addClass('btn-primary');
  $('#show-' + activeImage).removeClass('btn-default');
  $('#show-' + hiddenImage).addClass('btn-default');
  $('#show-' + hiddenImage).removeClass('btn-primary');
  $('#' + hiddenImage + '-image').hide();
  $('#' + activeImage + '-image').show();
}

$(function() {
  var currentOriginal = '/img/original.jpg';
  var currentResult = '/img/result.png';

  $('#sigma-slider').noUiSlider({
    start: 0.8,
    step: 0.05,
    connect: 'lower',
    range: {
      min: 0.0,
      max: 1.0
    }
  });

  $('#k-slider').noUiSlider({
    start: 300,
    step: 10,
    connect: 'lower',
    range: {
      min: 0,
      max: 3000
    }
  });

  $('#minsize-slider').noUiSlider({
    start: 50,
    step: 1,
    connect: 'lower',
    range: {
      min: 1,
      max: 1000
    }
  });

  $('#minweight-slider').noUiSlider({
    start: 10,
    step: 1,
    connect: 'lower',
    range: {
      min: 0,
      max: 50
    }
  });

  $('#initcredit-slider').noUiSlider({
    start: 100,
    step: 5,
    connect: 'lower',
    range: {
      min: 1,
      max: 1000
    }
  });

  $('#sigma-slider').Link('lower').to($('#input-sigma'));

  $('#k-slider').Link('lower').to($('#input-k'), null, {
    to: parseInt,
    from: Number
  });

  $('#minsize-slider').Link('lower').to($('#input-minsize'), null, {
    to: parseInt,
    from: Number
  });

  $('#minweight-slider').Link('lower').to($('#input-minweight'), null, {
    to: parseInt,
    from: Number
  });

  $('#initcredit-slider').Link('lower').to($('#input-initcredit'), null, {
    to: parseInt,
    from: Number
  });

  $('#algorithm-select').change(function() {
    if ($("#algorithm-select option:selected").val() == 1) {
      $('#phsmf-params').hide();
      $('#gbs-params').show();
    } else {
      $('#phsmf-params').show();
      $('#gbs-params').hide();
    }
  });

  $('#show-original').click(function() {
    changeImage('original', 'result');
  });

  $('#show-result').click(function() {
    changeImage('result', 'original');
  });

  $('#settings-form').submit(function() {
    $('#btn-run').attr('disabled', 'disabled');
    $('#btn-run').text('Segmenting');
    $('#work-status').show();
    $.ajax({
      type: 'POST',
      url: '/segment',
      data: new FormData(this),
      processData: false,
      contentType: false,
      success: function(data) {
        $('#btn-run').removeAttr('disabled');
        $('#btn-run').text('Run');
        $('#work-status').hide();
        data = data.split(' ');
        var filename = data[0];
        var originalext = data[1];
        current_original = filename + originalext;
        current_result  = 'new_' + filename + '.png';
        $('#original-image').attr('src', '/tmp/' + current_original);
        $('#result-image').attr('src', '/tmp/' + current_result);
        changeImage('result', 'original');
      }
    });
    return false;
  });

});
